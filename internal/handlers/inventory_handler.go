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
               i.quantity, i.reorder_point, i.max_stock, i.last_updated
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
        Quantity      int     `json:"quantity"`
        ReorderPoint  int     `json:"reorder_point"`
        MaxStock      int     `json:"max_stock"`
        NeedsReorder  bool    `json:"needs_reorder"`
    }

    var inventory []InventoryResponse
    for rows.Next() {
        var inv InventoryResponse
        err := rows.Scan(&inv.ID, &inv.ProductID, &inv.ProductName, &inv.ProductSKU,
            &inv.WarehouseID, &inv.WarehouseName, &inv.Quantity, &inv.ReorderPoint, &inv.MaxStock)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        inv.NeedsReorder = inv.Quantity <= inv.ReorderPoint
        inventory = append(inventory, inv)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(inventory)
}

// GetLowStock returns all products that need reordering (GET /api/inventory/low-stock)
func GetLowStock(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT i.product_id, p.name, i.warehouse_id, w.name, i.quantity, i.reorder_point,
               (i.reorder_point - i.quantity + 10) as recommended_qty
        FROM inventory i
        JOIN products p ON i.product_id = p.id
        JOIN warehouses w ON i.warehouse_id = w.id
        WHERE i.quantity <= i.reorder_point
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
            &item.WarehouseName, &item.CurrentStock, &item.ReorderPoint, &item.RecommendedQty)
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
    `, req.ProductID, req.WarehouseID, req.Quantity, req.Reason, req.PONumber, 1, currentQty)
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