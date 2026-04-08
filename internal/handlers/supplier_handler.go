package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetSuppliers returns all suppliers (GET /api/suppliers)
func GetSuppliers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, contact_person, phone, email, address, tax_id,
               payment_terms, lead_time_days, minimum_order, rating, 
               total_orders, on_time_deliveries, created_at, updated_at
        FROM suppliers ORDER BY name
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var suppliers []models.Supplier
    for rows.Next() {
        var s models.Supplier
        err := rows.Scan(
            &s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Email, &s.Address, &s.TaxID,
            &s.PaymentTerms, &s.LeadTimeDays, &s.MinimumOrder, &s.Rating,
            &s.TotalOrders, &s.OnTimeDeliveries, &s.CreatedAt, &s.UpdatedAt,
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        suppliers = append(suppliers, s)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(suppliers)
}

// GetSupplier returns a single supplier by ID (GET /api/suppliers/{id})
func GetSupplier(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
        return
    }

    var supplier models.Supplier
    err = db.DB.QueryRow(`
        SELECT id, name, contact_person, phone, email, address, tax_id,
               payment_terms, lead_time_days, minimum_order, rating, 
               total_orders, on_time_deliveries, created_at, updated_at
        FROM suppliers WHERE id = $1
    `, id).Scan(
        &supplier.ID, &supplier.Name, &supplier.ContactPerson, &supplier.Phone, &supplier.Email,
        &supplier.Address, &supplier.TaxID, &supplier.PaymentTerms, &supplier.LeadTimeDays,
        &supplier.MinimumOrder, &supplier.Rating, &supplier.TotalOrders,
        &supplier.OnTimeDeliveries, &supplier.CreatedAt, &supplier.UpdatedAt,
    )
    
    if err != nil {
        http.Error(w, "Supplier not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(supplier)
}
