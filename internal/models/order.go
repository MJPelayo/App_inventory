package models

import "time"

// SalesOrder represents a customer purchase order
type SalesOrder struct {
    ID              int       `json:"id"`                // Unique order identifier
    OrderNumber     string    `json:"order_number"`      // Human-readable order number (SO-2024-001)
    CustomerName    string    `json:"customer_name"`     // Name of the customer
    CustomerEmail   string    `json:"customer_email"`    // Customer email for notifications
    CustomerPhone   string    `json:"customer_phone"`    // Customer contact number
    ShippingAddress string    `json:"shipping_address"`  // Delivery address (if delivery)
    DeliveryType    string    `json:"delivery_type"`     // 'delivery' or 'pickup'
    Status          string    `json:"status"`            // pending/processing/completed/cancelled/shipped/delivered
    PaymentStatus   string    `json:"payment_status"`    // pending/partial/paid/refunded
    Subtotal        float64   `json:"subtotal"`          // Sum of all item totals before discounts
    DiscountAmount  float64   `json:"discount_amount"`   // Total discount applied to order
    Tax             float64   `json:"tax"`               // Tax amount (if applicable)
    ShippingCost    float64   `json:"shipping_cost"`     // Delivery fee
    TotalAmount     float64   `json:"total_amount"`      // Final amount after all adjustments
    CreatedBy       int       `json:"created_by"`        // User ID who created this order
    CreatedAt       time.Time `json:"created_at"`        // Order creation timestamp
    Items           []OrderItem `json:"items,omitempty"` // List of items in this order
}

// OrderItem represents a single product line within an order
type OrderItem struct {
    ID          int     `json:"id"`           // Unique item identifier
    OrderID     int     `json:"order_id"`     // Parent order ID
    OrderType   string  `json:"order_type"`   // 'sales' or 'supply'
    ProductID   int     `json:"product_id"`   // Which product was ordered
    ProductName string  `json:"product_name"` // Product name (snapshot at order time)
    Quantity    int     `json:"quantity"`     // Number of units ordered
    UnitPrice   float64 `json:"unit_price"`   // Price per unit at order time
    Discount    float64 `json:"discount"`     // Discount percentage for this item (0-100)
}

// CreateOrderRequest is used when creating a new sales order
type CreateOrderRequest struct {
    CustomerName    string             `json:"customer_name"`     // Required: customer's full name
    CustomerEmail   string             `json:"customer_email"`    // Optional: customer email
    CustomerPhone   string             `json:"customer_phone"`    // Optional: customer phone
    ShippingAddress string             `json:"shipping_address"`  // Required for delivery
    DeliveryType    string             `json:"delivery_type"`     // Required: 'delivery' or 'pickup'
    Items           []OrderItemRequest `json:"items"`             // Required: list of products
}

// OrderItemRequest represents a product being added to a new order
type OrderItemRequest struct {
    ProductID int `json:"product_id"` // Required: ID of product being ordered
    Quantity  int `json:"quantity"`   // Required: how many units (must be positive)
}

// SupplyOrder represents a purchase order to a supplier
type SupplyOrder struct {
    ID                  int        `json:"id"`                     // Unique PO identifier
    PONumber            string     `json:"po_number"`              // Purchase order number (PO-2024-001)
    SupplierID          int        `json:"supplier_id"`            // Supplier being ordered from
    SupplierName        string     `json:"supplier_name"`          // Supplier name (for display)
    Status              string     `json:"status"`                 // pending/processing/shipped/delivered/cancelled
    OrderDate           time.Time  `json:"order_date"`             // When order was created
    ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"` // When supplier should deliver
    ActualDeliveryDate  *time.Time `json:"actual_delivery_date"`   // When actually delivered
    Subtotal            float64    `json:"subtotal"`               // Sum of all item costs
    ShippingCost        float64    `json:"shipping_cost"`          // Delivery fee from supplier
    TotalAmount         float64    `json:"total_amount"`           // Subtotal + shipping
    CreatedBy           int        `json:"created_by"`             // User who created PO
    Notes               string     `json:"notes"`                  // Internal notes
    Items               []OrderItem `json:"items,omitempty"`       // Products being ordered
}