package models

import "time"

// StockMovement records every inventory change for complete audit trail
type StockMovement struct {
    ID              int       `json:"id"`                // Unique movement record ID
    ProductID       int       `json:"product_id"`        // Product that was moved
    WarehouseID     int       `json:"warehouse_id"`      // Warehouse where movement occurred
    QuantityChange  int       `json:"quantity_change"`   // Positive for additions, negative for removals
    MovementType    string    `json:"movement_type"`     // received/sold/transferred/adjusted/returned/damaged
    Reason          string    `json:"reason"`            // Human-readable reason for movement
    ReferenceNumber string    `json:"reference_number"`  // Order number, PO number, etc.
    PerformedBy     int       `json:"performed_by"`      // User ID who performed the action
    PreviousQuantity int      `json:"previous_quantity"` // Quantity before the movement
    NewQuantity     int       `json:"new_quantity"`      // Quantity after the movement
    CreatedAt       time.Time `json:"created_at"`        // When the movement occurred
}