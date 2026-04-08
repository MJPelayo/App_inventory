package models

import "time"

// Product represents catalog information (NO stock fields - stock is in Inventory)
// This separation allows one product to exist in multiple warehouses
type Product struct {
    ID          int       `json:"id"`           // Unique product identifier
    Name        string    `json:"name"`         // Product display name
    Description string    `json:"description"`  // Detailed product description
    SKU         string    `json:"sku"`          // Stock keeping unit (must be unique)
    Brand       string    `json:"brand"`        // Manufacturer brand name
    Price       float64   `json:"price"`        // Selling price to customers
    Cost        float64   `json:"cost"`         // Purchase cost from supplier
    CategoryID  int       `json:"category_id"`  // Foreign key to categories table
    SupplierID  int       `json:"supplier_id"`  // Foreign key to suppliers table
    ImageURL    *string   `json:"image_url"`    // Optional product image path (nullable)
    IsActive    bool      `json:"is_active"`    // Soft delete flag
    CreatedAt   time.Time `json:"created_at"`   // Record creation timestamp
    UpdatedAt   time.Time `json:"updated_at"`   // Last update timestamp
}

// ProductWithInventory combines product data with its stock levels across warehouses
type ProductWithInventory struct {
    Product   Product              `json:"product"`   // Product catalog details
    Inventory []InventoryWithLocation `json:"inventory"` // Stock levels with locations
}

// CreateProductRequest is used when adding a new product to catalog
type CreateProductRequest struct {
    Name        string  `json:"name"`        // Required: product name
    Description string  `json:"description"` // Optional: product description
    SKU         string  `json:"sku"`         // Required: unique SKU
    Brand       string  `json:"brand"`       // Optional: brand name
    Price       float64 `json:"price"`       // Required: selling price
    Cost        float64 `json:"cost"`        // Required: purchase cost
    CategoryID  int     `json:"category_id"` // Optional: category assignment
    SupplierID  int     `json:"supplier_id"` // Optional: supplier assignment
    ImageURL    string  `json:"image_url"`   // Optional: image URL
}

// UpdateProductRequest is used when editing an existing product
type UpdateProductRequest struct {
    Name        string  `json:"name"`        // Optional: new name
    Description string  `json:"description"` // Optional: new description
    Price       float64 `json:"price"`       // Optional: new price
    Cost        float64 `json:"cost"`        // Optional: new cost
    CategoryID  int     `json:"category_id"` // Optional: new category
    SupplierID  int     `json:"supplier_id"` // Optional: new supplier
    IsActive    *bool   `json:"is_active"`   // Optional: activate/deactivate
}