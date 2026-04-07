package models

import "time"

// Warehouse represents a physical storage location
type Warehouse struct {
    ID               int       `json:"id"`                 // Unique warehouse identifier
    Name             string    `json:"name"`               // Warehouse display name
    Location         string    `json:"location"`           // City/region where warehouse is located
    Capacity         int       `json:"capacity"`           // Maximum units this warehouse can hold
    CurrentOccupancy int       `json:"current_occupancy"`  // Current total units stored
    IsActive         bool      `json:"is_active"`          // Whether warehouse is operational
    CreatedAt        time.Time `json:"created_at"`         // Warehouse creation timestamp
    UpdatedAt        time.Time `json:"updated_at"`         // Last update timestamp
}

// CreateWarehouseRequest is used when adding a new warehouse
type CreateWarehouseRequest struct {
    Name     string `json:"name"`     // Required: warehouse name
    Location string `json:"location"` // Optional: location description
    Capacity int    `json:"capacity"` // Optional: maximum capacity (default 0 = unlimited)
}

// UpdateWarehouseRequest is used when editing warehouse details
type UpdateWarehouseRequest struct {
    Name     string `json:"name"`     // Optional: new name
    Location string `json:"location"` // Optional: new location
    Capacity *int   `json:"capacity"` // Optional: new capacity
    IsActive *bool  `json:"is_active"` // Optional: activate/deactivate
}