package handlers

import (
    "encoding/json"
    "net/http"

    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetInventory returns all inventory records with product and warehouse names (GET /api/inventory)
func GetInventory(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT i.id, i.product_id, p.name, p.sku, i.warehouse_id, w.name,
               i.quantity, i.reserved_quantity, i.reorder_point, i.max_stock, i.last_updated
        FROM inventory i
        JOIN products p ON i.product_id = p.id
        JOIN warehouses w ON i.warehouse_id = w.id
        ORDER BY i.id
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type InventoryResponse struct {
        ID            int     `json:"id"`
        ProductID     int     `json:"product_id"`
        ProductName   string  `json:"product_name"`
        ProductSKU    string  `json:"product_sku"`
        WarehouseID   int     `json:"warehouse_id"`
        WarehouseName string  `json:"warehouse_name"`
        Quantity        int     `json:"quantity"`
        ReservedQuantity int     `json:"reserved_quantity"`
        AvailableQuantity int    `json:"available_quantity"`
        ReorderPoint    int     `json:"reorder_point"`
        MaxStock        int     `json:"max_stock"`
        NeedsReorder    bool    `json:"needs_reorder"`
    }

    var inventory []InventoryResponse
    for rows.Next() {
        var inv InventoryResponse
        err := rows.Scan(&inv.ID, &inv.ProductID, &inv.ProductName, &inv.ProductSKU,
            &inv.WarehouseID, &inv.WarehouseName, &inv.Quantity, &inv.ReservedQuantity,
            &inv.ReorderPoint, &inv.MaxStock)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        inv.AvailableQuantity = inv.Quantity - inv.ReservedQuantity
        inv.NeedsReorder = inv.AvailableQuantity <= inv.ReorderPoint
        inventory = append(inventory, inv)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(inventory)
}

// GetLowStock returns all products that need reordering (GET /api/inventory/low-stock)
func GetLowStock(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT i.product_id, p.name, i.warehouse_id, w.name, i.quantity, i.reserved_quantity, i.reorder_point,
               (i.reorder_point - (i.quantity - i.reserved_quantity) + 10) as recommended_qty
        FROM inventory i
        JOIN products p ON i.product_id = p.id
        JOIN warehouses w ON i.warehouse_id = w.id
        WHERE (i.quantity - i.reserved_quantity) <= i.reorder_point
        ORDER BY (i.quantity::float / i.reorder_point) ASC
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var items []models.LowStockItem
    for rows.Next() {
        var item models.LowStockItem
        err := rows.Scan(&item.ProductID, &item.ProductName, &item.WarehouseID,
            &item.WarehouseName, &item.CurrentStock, &item.ReservedQuantity, &item.ReorderPoint, &item.RecommendedQty)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        items = append(items, item)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

// ReceiveStock processes incoming inventory (POST /api/inventory/receive)
func ReceiveStock(w http.ResponseWriter, r *http.Request) {
    var req models.ReceiveStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        http.Error(w, "Quantity must be positive", http.StatusBadRequest)
        return
    }

    // Start a transaction
    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Get current quantity before update
    var currentQty int
    err = tx.QueryRow(`
        SELECT quantity FROM inventory 
        WHERE product_id = $1 AND warehouse_id = $2
    `, req.ProductID, req.WarehouseID).Scan(&currentQty)
    if err != nil {
        http.Error(w, "Product not found in specified warehouse", http.StatusNotFound)
        return
    }

    // Update inventory quantity
    _, err = tx.Exec(`
        UPDATE inventory 
        SET quantity = quantity + $1, last_updated = CURRENT_TIMESTAMP
        WHERE product_id = $2 AND warehouse_id = $3
    `, req.Quantity, req.ProductID, req.WarehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Record stock movement for audit trail
    _, err = tx.Exec(`
        INSERT INTO stock_movements (product_id, warehouse_id, quantity_change, movement_type, 
                                     reason, reference_number, performed_by, previous_quantity, new_quantity)
        VALUES ($1, $2, $3, 'received', $4, $5, $6, $7, $7 + $3)
    `, req.ProductID, req.WarehouseID, req.Quantity, req.Reason, req.PONumber, 1, currentQty) // TODO: Replace 1 with actual user ID from JWT context
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Stock received successfully",
        "product_id": req.ProductID,
        "warehouse_id": req.WarehouseID,
        "quantity_received": req.Quantity,
    })
}

// ReserveInventory reserves stock for an order (POST /api/inventory/reserve)
func ReserveInventory(w http.ResponseWriter, r *http.Request) {
    var req models.ReserveStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        http.Error(w, "Quantity must be positive", http.StatusBadRequest)
        return
    }

    // Start a transaction
    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Check available quantity
    var availableQty int
    err = tx.QueryRow(`
        SELECT quantity - reserved_quantity
        FROM inventory
        WHERE product_id = $1 AND warehouse_id = $2
    `, req.ProductID, req.WarehouseID).Scan(&availableQty)
    if err != nil {
        http.Error(w, "Product not found in specified warehouse", http.StatusNotFound)
        return
    }

    if availableQty < req.Quantity {
        http.Error(w, "Insufficient available stock", http.StatusBadRequest)
        return
    }

    // Reserve the inventory
    _, err = tx.Exec(`
        UPDATE inventory
        SET reserved_quantity = reserved_quantity + $1, last_updated = CURRENT_TIMESTAMP
        WHERE product_id = $2 AND warehouse_id = $3
    `, req.Quantity, req.ProductID, req.WarehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Record stock movement for audit trail
    _, err = tx.Exec(`
        INSERT INTO stock_movements (product_id, warehouse_id, quantity_change, movement_type,
                                     reason, reference_number, performed_by, previous_quantity, new_quantity)
        VALUES ($1, $2, $3, 'reserved', $4, $5, $6, $7, $7)
    `, req.ProductID, req.WarehouseID, req.Quantity, "Order reservation", req.OrderID, 1, availableQty) // TODO: Replace 1 with actual user ID from JWT context
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Stock reserved successfully",
        "product_id": req.ProductID,
        "warehouse_id": req.WarehouseID,
        "quantity_reserved": req.Quantity,
        "order_id": req.OrderID,
    })
}

// ReleaseInventory releases reserved stock (POST /api/inventory/release)
func ReleaseInventory(w http.ResponseWriter, r *http.Request) {
    var req models.ReserveStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        http.Error(w, "Quantity must be positive", http.StatusBadRequest)
        return
    }

    // Start a transaction
    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Check current reserved quantity
    var currentReserved int
    err = tx.QueryRow(`
        SELECT reserved_quantity
        FROM inventory
        WHERE product_id = $1 AND warehouse_id = $2
    `, req.ProductID, req.WarehouseID).Scan(&currentReserved)
    if err != nil {
        http.Error(w, "Product not found in specified warehouse", http.StatusNotFound)
        return
    }

    if currentReserved < req.Quantity {
        http.Error(w, "Cannot release more than reserved quantity", http.StatusBadRequest)
        return
    }

    // Release the inventory
    _, err = tx.Exec(`
        UPDATE inventory
        SET reserved_quantity = reserved_quantity - $1, last_updated = CURRENT_TIMESTAMP
        WHERE product_id = $2 AND warehouse_id = $3
    `, req.Quantity, req.ProductID, req.WarehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Record stock movement for audit trail
    _, err = tx.Exec(`
        INSERT INTO stock_movements (product_id, warehouse_id, quantity_change, movement_type,
                                     reason, reference_number, performed_by, previous_quantity, new_quantity)
        VALUES ($1, $2, $3, 'released', $4, $5, $6, $7, $7)
    `, req.ProductID, req.WarehouseID, req.Quantity, "Order cancellation/release", req.OrderID, 1, currentReserved) // TODO: Replace 1 with actual user ID from JWT context
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Stock released successfully",
        "product_id": req.ProductID,
        "warehouse_id": req.WarehouseID,
        "quantity_released": req.Quantity,
        "order_id": req.OrderID,
    })
}