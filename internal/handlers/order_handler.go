package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    "inventory-api/internal/db"
    "inventory-api/internal/models"
)

// CreateSalesOrder creates a new customer sales order (POST /api/sales-orders)
func CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
    var req models.CreateOrderRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if len(req.Items) == 0 {
        http.Error(w, "Order must have at least one item", http.StatusBadRequest)
        return
    }

    // Start transaction
    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Generate order number
    orderNumber := "SO-" + strconv.FormatInt(time.Now().Unix(), 10)

    // Calculate subtotal
    var subtotal float64
    for _, item := range req.Items {
        var price float64
        err := tx.QueryRow("SELECT price FROM products WHERE id = $1", item.ProductID).Scan(&price)
        if err != nil {
            http.Error(w, "Product not found", http.StatusBadRequest)
            return
        }
        subtotal += price * float64(item.Quantity)
    }

    // Create the order
    var orderID int
    err = tx.QueryRow(`
        INSERT INTO sales_orders (order_number, customer_name, customer_email, customer_phone, 
                                  shipping_address, delivery_type, subtotal, total_amount, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $8)
        RETURNING id
    `, orderNumber, req.CustomerName, req.CustomerEmail, req.CustomerPhone,
        req.ShippingAddress, req.DeliveryType, subtotal, 1).Scan(&orderID) // TODO: Replace 1 with actual user ID from JWT context
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Add order items
    for _, item := range req.Items {
        var productName string
        var price float64
        tx.QueryRow("SELECT name, price FROM products WHERE id = $1", item.ProductID).Scan(&productName, &price)

        _, err = tx.Exec(`
            INSERT INTO order_items (order_id, order_type, product_id, product_name, quantity, unit_price)
            VALUES ($1, 'sales', $2, $3, $4, $5)
        `, orderID, item.ProductID, productName, item.Quantity, price)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "order_id": orderID,
        "order_number": orderNumber,
        "status": "pending",
        "message": "Sales order created successfully",
    })
}

// GetSalesOrders returns all sales orders (GET /api/sales-orders)
func GetSalesOrders(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query(`
        SELECT id, order_number, customer_name, customer_email, customer_phone,
               shipping_address, delivery_type, status, payment_status,
               subtotal, discount_amount, tax, shipping_cost, total_amount,
               created_by, created_at
        FROM sales_orders ORDER BY created_at DESC
    `)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var orders []models.SalesOrder
    for rows.Next() {
        var o models.SalesOrder
        err := rows.Scan(
            &o.ID, &o.OrderNumber, &o.CustomerName, &o.CustomerEmail, &o.CustomerPhone,
            &o.ShippingAddress, &o.DeliveryType, &o.Status, &o.PaymentStatus,
            &o.Subtotal, &o.DiscountAmount, &o.Tax, &o.ShippingCost, &o.TotalAmount,
            &o.CreatedBy, &o.CreatedAt,
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        orders = append(orders, o)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(orders)
}

// GetSalesOrder returns a single sales order with its items (GET /api/sales-orders/{id})
func GetSalesOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    var order models.SalesOrder
    err = db.DB.QueryRow(`
        SELECT id, order_number, customer_name, customer_email, customer_phone,
               shipping_address, delivery_type, status, payment_status,
               subtotal, discount_amount, tax, shipping_cost, total_amount,
               created_by, created_at
        FROM sales_orders WHERE id = $1
    `, id).Scan(
        &order.ID, &order.OrderNumber, &order.CustomerName, &order.CustomerEmail, &order.CustomerPhone,
        &order.ShippingAddress, &order.DeliveryType, &order.Status, &order.PaymentStatus,
        &order.Subtotal, &order.DiscountAmount, &order.Tax, &order.ShippingCost, &order.TotalAmount,
        &order.CreatedBy, &order.CreatedAt,
    )
    
    if err != nil {
        http.Error(w, "Order not found", http.StatusNotFound)
        return
    }

    // Get order items
    rows, err := db.DB.Query(`
        SELECT id, product_id, product_name, quantity, unit_price, discount
        FROM order_items WHERE order_id = $1 AND order_type = 'sales'
    `, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var items []models.OrderItem
    for rows.Next() {
        var item models.OrderItem
        err := rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPrice, &item.Discount)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        items = append(items, item)
    }
    order.Items = items

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}
