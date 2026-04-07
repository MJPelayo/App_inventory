-- Individual items within an order (works for both sales and supply orders)
-- order_type distinguishes between 'sales' and 'supply' orders
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    order_type VARCHAR(10) NOT NULL CHECK (order_type IN ('sales', 'supply')),
    product_id INTEGER NOT NULL,
    product_name VARCHAR(200),  -- Snapshot of product name at order time
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    discount DECIMAL(5,2) DEFAULT 0 CHECK (discount >= 0 AND discount <= 100),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- Indexes for order item queries
CREATE INDEX idx_order_items_order ON order_items(order_id, order_type);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);