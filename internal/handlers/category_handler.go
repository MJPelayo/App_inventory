package handlers

import (
    "encoding/json"
    "net/http"

    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// GetCategories returns all categories (GET /api/categories)
func GetCategories(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, name, parent_id, description, created_at
        FROM categories ORDER BY parent_id, name
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var c models.Category
        err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.Description, &c.CreatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        categories = append(categories, c)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

// GetCategoryTree returns hierarchical category tree (GET /api/categories/tree)
func GetCategoryTree(w http.ResponseWriter, r *http.Request) {
    // Get all categories first
    rows, err := db.DB.Query(`
        SELECT id, name, parent_id, description
        FROM categories ORDER BY id
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Build map of categories
    categoryMap := make(map[int]*models.CategoryTree)
    var rootCategories []models.CategoryTree

    for rows.Next() {
        var id, parentID int
        var name, description string
        rows.Scan(&id, &name, &parentID, &description)
        
        cat := models.CategoryTree{
            ID:          id,
            Name:        name,
            Description: description,
            Children:    []models.CategoryTree{},
        }
        categoryMap[id] = &cat

        if parentID == 0 {
            rootCategories = append(rootCategories, cat)
        } else {
            if parent, exists := categoryMap[parentID]; exists {
                parent.Children = append(parent.Children, cat)
            }
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rootCategories)
}