package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetWarehouses returns all warehouses (GET /api/warehouses)
func GetWarehouses(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, location, capacity, current_occupancy, is_active, created_at, updated_at
        FROM warehouses ORDER BY id
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var warehouses []models.Warehouse
    for rows.Next() {
        var warehouse models.Warehouse
        err := rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Location, &warehouse.Capacity, &warehouse.CurrentOccupancy, &warehouse.IsActive, &warehouse.CreatedAt, &warehouse.UpdatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        warehouses = append(warehouses, warehouse)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(warehouses)
}

// GetProductLocations returns all product locations for a warehouse (GET /api/warehouses/{id}/locations)
func GetProductLocations(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    warehouseID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
        return
    }

    rows, err := db.DB.Query(`
        SELECT pl.id, pl.product_id, p.name, p.sku, pl.aisle_number, pl.side, 
               pl.shelf_number, pl.layer, pl.quantity
        FROM product_locations pl
        JOIN products p ON pl.product_id = p.id
        WHERE pl.warehouse_id = $1
        ORDER BY pl.aisle_number, pl.side, pl.shelf_number
    `, warehouseID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type LocationResponse struct {
        ID          int    `json:"id"`
        ProductID   int    `json:"product_id"`
        ProductName string `json:"product_name"`
        ProductSKU  string `json:"product_sku"`
        Aisle       int    `json:"aisle"`
        Side        string `json:"side"`
        Shelf       int    `json:"shelf"`
        Layer       string `json:"layer"`
        Quantity    int    `json:"quantity"`
    }

    var locations []LocationResponse
    for rows.Next() {
        var loc LocationResponse
        err := rows.Scan(&loc.ID, &loc.ProductID, &loc.ProductName, &loc.ProductSKU,
            &loc.Aisle, &loc.Side, &loc.Shelf, &loc.Layer, &loc.Quantity)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        locations = append(locations, loc)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(locations)
}