// Package main is the entry point for the Inventory Management System API
// This file starts the HTTP server and initializes database and routes
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv" // Loads environment variables from .env file
    
    "inventory-api/internal/db"
    "inventory-api/internal/routes"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Get database path from environment or use default
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./inventory.db"
    }

    // Initialize database connection and run migrations
    if err := db.InitDB(dbPath); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.CloseDB() // Close connection when program exits

    // Get port from environment or use default 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Setup all API routes
    router := routes.SetupRoutes()

    // Start the HTTP server
    log.Printf("🚀 Server starting on http://localhost:%s", port)
    log.Printf("📋 Test the API with:")
    log.Printf("   GET  http://localhost:%s/api/products", port)
    log.Printf("   GET  http://localhost:%s/api/inventory", port)
    log.Printf("   POST http://localhost:%s/api/auth/login", port)
    
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}