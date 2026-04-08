package models

import "time"

// Inventory tracks stock levels for a product in a specific warehouse
// This is the SINGLE SOURCE OF TRUTH for all stock levels
type Inventory struct {
	ID               int       `json:"id"`                // Unique inventory record ID
	ProductID        int       `json:"product_id"`        // Which product this stock is for
	WarehouseID      int       `json:"warehouse_id"`      // Which warehouse holds this stock
	Quantity         int       `json:"quantity"`          // Current stock count (cannot be negative)
	ReservedQuantity int       `json:"reserved_quantity"` // Quantity reserved for pending orders
	ReorderPoint     int       `json:"reorder_point"`     // Threshold that triggers reorder alert
	MaxStock         int       `json:"max_stock"`         // Maximum allowed stock (overstock alert)
	LastUpdated      time.Time `json:"last_updated"`      // When stock was last modified
}

// InventoryWithLocation includes the physical storage location of inventory
type InventoryWithLocation struct {
	Inventory
	Locations []ProductLocation `json:"locations"` // Physical locations within warehouse
}

// ProductLocation tracks exact physical location using aisle/side/shelf/layer system
// No bins - uses human-readable location format
type ProductLocation struct {
	ID          int    `json:"id"`           // Unique location record ID
	ProductID   int    `json:"product_id"`   // Product stored at this location
	WarehouseID int    `json:"warehouse_id"` // Warehouse containing this location
	AisleNumber int    `json:"aisle_number"` // Aisle number (1, 2, 3, etc.)
	Side        string `json:"side"`         // 'left' or 'right' side of aisle
	ShelfNumber int    `json:"shelf_number"` // Shelf number within the aisle
	Layer       string `json:"layer"`        // 'top', 'middle', 'middle2', 'middle3', 'bottom'
	Quantity    int    `json:"quantity"`     // How many units at this location
}

// ReceiveStockRequest is used when warehouse receives new inventory
type ReceiveStockRequest struct {
	ProductID   int    `json:"product_id"`   // Required: which product is being received
	WarehouseID int    `json:"warehouse_id"` // Required: which warehouse is receiving
	Quantity    int    `json:"quantity"`     // Required: number of units received (must be positive)
	Reason      string `json:"reason"`       // Optional: why receiving (PO number, return, etc.)
	PONumber    string `json:"po_number"`    // Optional: purchase order reference
}

// TransferStockRequest is used when moving stock between warehouses
type TransferStockRequest struct {
	ProductID       int    `json:"product_id"`        // Required: which product to transfer
	FromWarehouseID int    `json:"from_warehouse_id"` // Required: source warehouse
	ToWarehouseID   int    `json:"to_warehouse_id"`   // Required: destination warehouse
	Quantity        int    `json:"quantity"`          // Required: how many units to transfer
	Reason          string `json:"reason"`            // Optional: why transferring
}

// ReserveStockRequest is used when creating an order to reserve inventory
type ReserveStockRequest struct {
    ProductID   int `json:"product_id"`    // Required: which product to reserve
    WarehouseID int `json:"warehouse_id"`  // Required: which warehouse to reserve from
    Quantity    int `json:"quantity"`      // Required: how many units to reserve
    OrderID     int `json:"order_id"`      // Required: order ID for reference
}

// LowStockItem represents a product that needs reordering
type LowStockItem struct {
	ProductID        int    `json:"product_id"`        // Product identifier
	ProductName      string `json:"product_name"`      // Product display name
	WarehouseID      int    `json:"warehouse_id"`      // Warehouse where stock is low
	WarehouseName    string `json:"warehouse_name"`    // Warehouse display name
	CurrentStock     int    `json:"current_stock"`     // Current quantity available
	ReservedQuantity int    `json:"reserved_quantity"` // Quantity reserved for orders
	ReorderPoint     int    `json:"reorder_point"`     // Threshold that triggered alert
	RecommendedQty   int    `json:"recommended_qty"`   // Suggested reorder quantity
}
