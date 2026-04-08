package handlers

import (
    "encoding/json"
    "net/http"

    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// ReserveStock reserves inventory for an order (POST /api/inventory/reserve)
func ReserveStock(w http.ResponseWriter, r *http.Request) {
    var req models.ReserveStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Quantity <= 0 {
        http.Error(w, "Quantity must be positive", http.StatusBadRequest)
        return
    }

    // Start transaction
    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Check available stock (quantity - reserved_quantity)
    var availableQty int
    err = tx.QueryRow(`
        SELECT (quantity - reserved_quantity) as available
        FROM inventory 
        WHERE product_id = $1 AND warehouse_id = $2
        FOR UPDATE  -- Lock the row to prevent race conditions
    `, req.ProductID, req.WarehouseID).Scan(&availableQty)
    
    if err != nil {
        http.Error(w, "Product not found in specified warehouse", http.StatusNotFound)
        return
    }

    if availableQty < req.Quantity {
        http.Error(w, "Insufficient available stock", http.StatusConflict)
        return
    }

    // Reserve the stock
    _, err = tx.Exec(`
        UPDATE inventory 
        SET reserved_quantity = reserved_quantity + $1,
            last_updated = CURRENT_TIMESTAMP
        WHERE product_id = $2 AND warehouse_id = $3
    `, req.Quantity, req.ProductID, req.WarehouseID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Record reservation in audit trail
    _, err = tx.Exec(`
        INSERT INTO stock_movements (product_id, warehouse_id, quantity_change, movement_type, 
                                     reason, reference_number, performed_by, previous_quantity, new_quantity)
        VALUES ($1, $2, $3, 'reserved', $4, $5, $6, 
                (SELECT quantity FROM inventory WHERE product_id = $1 AND warehouse_id = $2),
                (SELECT quantity FROM inventory WHERE product_id = $1 AND warehouse_id = $2))
    `, req.ProductID, req.WarehouseID, -req.Quantity, 
       "Order reservation", req.OrderID, 1)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Stock reserved successfully",
    })
}

// ReleaseReservedStock releases reserved inventory (when order is cancelled or completed)
func ReleaseReservedStock(w http.ResponseWriter, r *http.Request) {
    var req models.ReserveStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Release the reserved stock
    result, err := tx.Exec(`
        UPDATE inventory 
        SET reserved_quantity = GREATEST(reserved_quantity - $1, 0),
            last_updated = CURRENT_TIMESTAMP极
        WHERE product_id = $2 AND warehouse_id = $3
          AND reserved_quantity >= $1
    `, req.Quantity, req.ProductID, req.WarehouseID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Not enough reserved stock to release", http.StatusConflict)
        return
    }

    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Reserved stock released successfully",
    })
}