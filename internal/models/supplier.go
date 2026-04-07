package models

import "time"

// Supplier represents a company that provides products to the business
type Supplier struct {
    ID               int       `json:"id"`                  // Unique supplier identifier
    Name             string    `json:"name"`                // Supplier company name
    ContactPerson    string    `json:"contact_person"`      // Primary contact name
    Phone            string    `json:"phone"`               // Contact phone number
    Email            string    `json:"email"`               // Contact email address
    Address          string    `json:"address"`             // Physical/billing address
    TaxID            string    `json:"tax_id"`              // Tax identification number
    PaymentTerms     string    `json:"payment_terms"`       // e.g., "Net 30", "Net 60"
    LeadTimeDays     int       `json:"lead_time_days"`      // Average delivery days
    MinimumOrder     int       `json:"minimum_order"`       // Minimum quantity required
    Rating           float64   `json:"rating"`              // Performance rating (0-5)
    TotalOrders      int       `json:"total_orders"`        // Total purchase orders placed
    OnTimeDeliveries int       `json:"on_time_deliveries"`  // Count of on-time deliveries
    CreatedAt        time.Time `json:"created_at"`          // Record creation timestamp
    UpdatedAt        time.Time `json:"updated_at"`          // Last update timestamp
}

// SupplierPerformance tracks supplier metrics over time
type SupplierPerformance struct {
    SupplierID          int     `json:"supplier_id"`
    SupplierName        string  `json:"supplier_name"`
    OnTimeRate          float64 `json:"on_time_rate"`        // Percentage of on-time deliveries
    AverageLeadTime     float64 `json:"average_lead_time"`   // Average days from order to delivery
    TotalOrders         int     `json:"total_orders"`
    TotalSpent          float64 `json:"total_spent"`
    QualityRating       float64 `json:"quality_rating"`      // Based on return rates
}

// CreateSupplierRequest is used when adding a new supplier
type CreateSupplierRequest struct {
    Name          string `json:"name"`           // Required: supplier name
    ContactPerson string `json:"contact_person"` // Optional: contact name
    Phone         string `json:"phone"`          // Optional: phone number
    Email         string `json:"email"`          // Optional: email address
    Address       string `json:"address"`        // Optional: address
    TaxID         string `json:"tax_id"`         // Optional: tax ID
    PaymentTerms  string `json:"payment_terms"`  // Optional: payment terms
    LeadTimeDays  int    `json:"lead_time_days"` // Optional: lead time (default 7)
    MinimumOrder  int    `json:"minimum_order"`  // Optional: minimum order (default 0)
}