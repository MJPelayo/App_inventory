-- Supplier purchase orders (what the company buys from suppliers)
-- status: pending, processing, shipped, delivered, cancelled
CREATE TABLE IF NOT EXISTS supply_orders (
    id SERIAL PRIMARY KEY,
    po_number VARCHAR(50) UNIQUE NOT NULL,
    supplier_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expected_delivery_date TIMESTAMP,
    actual_delivery_date TIMESTAMP,
    subtotal DECIMAL(12,2) DEFAULT 0 CHECK (subtotal >= 0),
    shipping_cost DECIMAL(12,2) DEFAULT 0 CHECK (shipping_cost >= 0),
    total_amount DECIMAL(12,2) DEFAULT 0 CHECK (total_amount >= 0),
    created_by INTEGER,
    notes TEXT,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for supply order queries
CREATE INDEX idx_supply_orders_po_number ON supply_orders(po_number);
CREATE INDEX idx_supply_orders_supplier_id ON supply_orders(supplier_id);
CREATE INDEX idx_supply_orders_status ON supply_orders(status);