package handlers

import (
    "encoding/json"
    "net/http"
)

// GetSalesReport returns sales analytics (GET /api/reports/sales)
// For Test 1, returns mock data (will implement full reporting in final project)
func GetSalesReport(w http.ResponseWriter, r *http.Request) {
    report := map[string]interface{}{
        "status": "partial",
        "message": "Full sales report coming in final project",
        "sample_data": map[string]interface{}{
            "total_sales": 3099.97,
            "total_orders": 1,
            "average_order_value": 3099.97,
            "top_product": "Galaxy S23",
            "period": "last_30_days",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(report)
}

// GetInventoryReport returns inventory analytics (GET /api/reports/inventory)
func GetInventoryReport(w http.ResponseWriter, r *http.Request) {
    report := map[string]interface{}{
        "status": "partial",
        "message": "Full inventory report coming in final project",
        "sample_data": map[string]interface{}{
            "total_products": 8,
            "total_units": 151,
            "low_stock_items": 2,
            "inventory_value": 45000.00,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(report)
}