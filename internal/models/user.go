package models
// Package models defines data structures for API requests and responses
package models

import "time"

// User represents a system user account with role-based access
type User struct {
    ID             int        `json:"id"`                       // Unique user identifier
    Name           string     `json:"name"`                     // User's full name
    Email          string     `json:"email"`                    // Login email (must be unique)
    PasswordHash   string     `json:"-"`                        // Hashed password (excluded from JSON responses)
    Role           string     `json:"role"`                     // admin/sales_manager/warehouse_manager/supply_manager
    Department     *string    `json:"department,omitempty"`     // Admin only: department name
    SalesTarget    *float64   `json:"sales_target,omitempty"`   // Sales manager: monthly sales goal
    CommissionRate *float64   `json:"commission_rate,omitempty"` // Sales manager: commission percentage (default 5%)
    WarehouseID    *int       `json:"warehouse_id,omitempty"`   // Warehouse manager: assigned warehouse
    Shift          *string    `json:"shift,omitempty"`          // Warehouse manager: Day/Night shift
    PurchaseBudget *float64   `json:"purchase_budget,omitempty"` // Supply manager: annual purchase budget
    CreatedAt      time.Time  `json:"created_at"`               // Account creation timestamp
    LastLogin      *time.Time `json:"last_login,omitempty"`     // Most recent login time
}

// LoginRequest contains credentials for user authentication
type LoginRequest struct {
    Email    string `json:"email"`    // User's email address
    Password string `json:"password"` // User's plain text password
}

// LoginResponse returns user info and auth token after successful login
type LoginResponse struct {
    Token string `json:"token"` // JWT authentication token (simple for Test 1)
    User  User   `json:"user"`  // User information (excluding password)
}

// CreateUserRequest is used when admin creates a new user account
type CreateUserRequest struct {
    Name           string   `json:"name"`                      // Required: full name
    Email          string   `json:"email"`                     // Required: unique email
    Password       string   `json:"password"`                  // Required: plain password (will be hashed)
    Role           string   `json:"role"`                      // Required: user role
    Department     *string  `json:"department,omitempty"`      // Optional: for admin users
    SalesTarget    *float64 `json:"sales_target,omitempty"`    // Optional: for sales managers
    CommissionRate *float64 `json:"commission_rate,omitempty"` // Optional: for sales managers
    WarehouseID    *int     `json:"warehouse_id,omitempty"`    // Optional: for warehouse managers
    Shift          *string  `json:"shift,omitempty"`           // Optional: for warehouse managers
    PurchaseBudget *float64 `json:"purchase_budget,omitempty"` // Optional: for supply managers
}