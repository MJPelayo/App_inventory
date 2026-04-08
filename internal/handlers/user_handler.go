package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetUsers returns all users from the database (GET /api/users)
// This endpoint is admin-only in production (simplified for Test 1)
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // Query all users from database
    rows, err := db.DB.Query(`
        SELECT id, name, email, role, department, sales_target, 
               commission_rate, warehouse_id, shift, purchase_budget, created_at
        FROM users ORDER BY id
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Iterate through results and build user list
    var users []models.User
    for rows.Next() {
        var u models.User
        err := rows.Scan(
            &u.ID, &u.Name, &u.Email, &u.Role, &u.Department,
            &u.SalesTarget, &u.CommissionRate, &u.WarehouseID,
            &u.Shift, &u.PurchaseBudget, &u.CreatedAt,
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// GetUser returns a single user by ID (GET /api/users/{id})
func GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user models.User
    err = db.DB.QueryRow(`
        SELECT id, name, email, role, department, sales_target, 
               commission_rate, warehouse_id, shift, purchase_budget, created_at
        FROM users WHERE id = $1
    `, id).Scan(
        &user.ID, &user.Name, &user.Email, &user.Role, &user.Department,
        &user.SalesTarget, &user.CommissionRate, &user.WarehouseID,
        &user.Shift, &user.PurchaseBudget, &user.CreatedAt,
    )
    
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// CreateUser adds a new user to the system (POST /api/users)
// Password is stored as plain text for Test 1 (will hash in final project)
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Insert new user into database
    var userID int
    err := db.DB.QueryRow(`
        INSERT INTO users (name, email, password_hash, role, department, 
                          sales_target, commission_rate, warehouse_id, shift, purchase_budget)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `, req.Name, req.Email, req.Password+"_hash", req.Role, req.Department,
       req.SalesTarget, req.CommissionRate, req.WarehouseID, req.Shift, req.PurchaseBudget,
    ).Scan(&userID)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the created user
    user, _ := GetUserByID(userID)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Helper function to get user by ID (used internally)
func GetUserByID(id int) (*models.User, error) {
    var user models.User
    err := db.DB.QueryRow(`
        SELECT id, name, email, role, department, created_at
        FROM users WHERE id = $1
    `, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Department, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}