// Package routes defines all API endpoint routes and their handlers
package routes

import (
    "net/http"

    "github.com/gorilla/mux"
    "inventory-api/internal/handlers"
)

// SetupRoutes configures all API routes and returns the router
func SetupRoutes() *mux.Router {
    router := mux.NewRouter()

    // API v1 base path
    api := router.PathPrefix("/api").Subrouter()

    // ==================== AUTH ROUTES ====================
    api.HandleFunc("/auth/login", handlers.Login).Methods("POST")

    // ==================== USER ROUTES ====================
    api.HandleFunc("/users", handlers.GetUsers).Methods("GET")
    api.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
    api.HandleFunc("/users", handlers.CreateUser).Methods("POST")

    // ==================== CATEGORY ROUTES ====================
    api.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
    api.HandleFunc("/categories/tree", handlers.GetCategoryTree).Methods("GET")

    // ==================== SUPPLIER ROUTES ====================
    api.HandleFunc("/suppliers", handlers.GetSuppliers).Methods("GET")
    api.HandleFunc("/suppliers/{id}", handlers.GetSupplier).Methods("GET")

    // ==================== PRODUCT ROUTES ====================
    api.HandleFunc("/products", handlers.GetProducts).Methods("GET")
    api.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
    api.HandleFunc("/products", handlers.CreateProduct).Methods("POST")

    // ==================== WAREHOUSE ROUTES ====================
    api.HandleFunc("/warehouses", handlers.GetWarehouses).Methods("GET")
    api.HandleFunc("/warehouses/{id}/locations", handlers.GetProductLocations).Methods("GET")

    // ==================== INVENTORY ROUTES ====================
    api.HandleFunc("/inventory", handlers.GetInventory).Methods("GET")
    api.HandleFunc("/inventory/low-stock", handlers.GetLowStock).Methods("GET")
    api.HandleFunc("/inventory/receive", handlers.ReceiveStock).Methods("POST")
    api.HandleFunc("/inventory/reserve", handlers.ReserveInventory).Methods("POST")
    api.HandleFunc("/inventory/release", handlers.ReleaseInventory).Methods("POST")

    // ==================== SALES ORDER ROUTES ====================
    api.HandleFunc("/sales-orders", handlers.GetSalesOrders).Methods("GET")
    api.HandleFunc("/sales-orders/{id}", handlers.GetSalesOrder).Methods("GET")
    api.HandleFunc("/sales-orders", handlers.CreateSalesOrder).Methods("POST")

    // ==================== REPORT ROUTES ====================
    api.HandleFunc("/reports/sales", handlers.GetSalesReport).Methods("GET")
    api.HandleFunc("/reports/inventory", handlers.GetInventoryReport).Methods("GET")

    // Health check endpoint
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    }).Methods("GET")

    return router
}
