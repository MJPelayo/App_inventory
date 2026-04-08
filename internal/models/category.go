package models

import "time"

// Category represents a product category with hierarchical parent-child relationship
type Category struct {
    ID          int       `json:"id"`           // Unique category identifier
    Name        string    `json:"name"`         // Category display name
    ParentID    int       `json:"parent_id"`    // 0 for root categories, otherwise parent category ID
    Description string    `json:"description"`  // Optional category description
    CreatedAt   time.Time `json:"created_at"`   // Creation timestamp
    Children    []Category `json:"children,omitempty"` // Child categories (for tree view)
}

// CategoryTree represents the hierarchical category structure
type CategoryTree struct {
    ID          int             `json:"id"`
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Children    []*CategoryTree `json:"children"`
}

// CreateCategoryRequest is used when adding a new category
type CreateCategoryRequest struct {
    Name        string `json:"name"`         // Required: category name
    ParentID    int    `json:"parent_id"`    // Optional: parent category (0 for root)
    Description string `json:"description"`  // Optional: category description
}