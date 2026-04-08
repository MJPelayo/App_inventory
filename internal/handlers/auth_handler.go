// Package handlers contains HTTP request handlers for each API endpoint
package handlers

import (
    "encoding/json"
    "net/http"

    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// Login authenticates a user and returns a session token (POST /api/auth/login)
// For Test 1, this uses simple password comparison (will upgrade to bcrypt + JWT later)
func Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    
    // Parse JSON request body into LoginRequest struct
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Query database for user with matching email
    var user models.User
    err := db.DB.QueryRow(`
        SELECT id, name, email, password_hash, role, department, 
               sales_target, commission_rate, warehouse_id, shift, 
               purchase_budget, created_at, last_login
        FROM users WHERE email = $1
    `, req.Email).Scan(
        &user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role,
        &user.Department, &user.SalesTarget, &user.CommissionRate, &user.WarehouseID,
        &user.Shift, &user.PurchaseBudget, &user.CreatedAt, &user.LastLogin,
    )
    
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Simple password check (Test 1 only - will implement bcrypt for final project)
    if user.PasswordHash != req.Password+"_hash" && user.PasswordHash != req.Password {
        // For sample data, passwords are stored as plain text with "_hash" suffix
        // In production, use bcrypt.CompareHashAndPassword()
        if user.PasswordHash != req.Password {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }
    }

    // Update last login timestamp
    db.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", user.ID)

    // For Test 1, return a simple token (will implement JWT for final project)
    response := models.LoginResponse{
        Token: "test-token-for-" + user.Email,
        User:  user,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
