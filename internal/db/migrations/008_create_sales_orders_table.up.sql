-- Customer sales orders with status and payment tracking
-- status: pending, processing, completed, cancelled, shipped, delivered
-- payment_status: pending, partial, paid, refunded
CREATE TABLE IF NOT EXISTS sales_orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_name VARCHAR(100) NOT NULL,
    customer_email VARCHAR(100),
    customer_phone VARCHAR(20),
    shipping_address TEXT,
    delivery_type VARCHAR(20) DEFAULT 'delivery' CHECK (delivery_type IN ('delivery', 'pickup')),
    status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'pending',
    subtotal DECIMAL(12,2) DEFAULT 0 CHECK (subtotal >= 0),
    discount_amount DECIMAL(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
    discount_approved_by INTEGER,
    tax DECIMAL(12,2) DEFAULT 0 CHECK (tax >= 0),
    shipping_cost DECIMAL(12,2) DEFAULT 0 CHECK (shipping_cost >= 0),
    total_amount DECIMAL(12,2) DEFAULT 0 CHECK (total_amount >= 0),
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (discount_approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for order queries
CREATE INDEX idx_sales_orders_order_number ON sales_orders(order_number);
CREATE INDEX idx_sales_orders_customer_name ON sales_orders(customer_name);
CREATE INDEX idx_sales_orders_status ON sales_orders(status);
CREATE INDEX idx_sales_orders_created_at ON sales_orders(created_at);