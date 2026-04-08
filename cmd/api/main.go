// Package main is the entry point for the Inventory Management System API
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    
    "inventory-api/internal/db"
    "inventory-api/internal/middleware"
    "inventory-api/internal/routes"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Initialize PostgreSQL database connection
    if err := db.InitDB(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.CloseDB()

    // Run database migrations (creates tables and seeds data)
    if err := db.RunMigrations(); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Setup all API routes
    router := routes.SetupRoutes()

    // Apply middleware
    router.Use(middleware.CORSMiddleware)
    router.Use(middleware.LoggingMiddleware)

    // Get port from environment or use default
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Inventory Management System API starting on http://localhost:%s", port)
    log.Printf("📋 Available endpoints:")
    log.Printf("   POST   /api/auth/login")
    log.Printf("   GET    /api/users")
    log.Printf("   GET    /api/products")
    log.Printf("   GET    /api/inventory")
    log.Printf("   GET    /api/inventory/low-stock")
    log.Printf("   POST   /api/inventory/receive")
    log.Printf("   GET    /api/sales-orders")
    log.Printf("   POST   /api/sales-orders")
    log.Printf("   GET    /api/categories/tree")
    log.Printf("   GET    /api/warehouses")
    log.Printf("   GET    /health")

    // Start the HTTP server
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}