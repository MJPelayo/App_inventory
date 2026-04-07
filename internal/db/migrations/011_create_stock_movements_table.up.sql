-- Complete audit trail for every stock change
-- Every inventory operation creates a record here
CREATE TABLE IF NOT EXISTS stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    quantity_change INTEGER NOT NULL,  -- Positive for additions, negative for removals
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('received', 'sold', 'transferred', 'adjusted', 'returned', 'damaged')),
    reason TEXT,
    reference_number VARCHAR(100),  -- Order number, PO number, etc.
    performed_by INTEGER,
    previous_quantity INTEGER NOT NULL,
    new_quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
    FOREIGN KEY (performed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for movement history queries
CREATE INDEX idx_stock_movements_product ON stock_movements(product_id);
CREATE INDEX idx_stock_movements_warehouse ON stock_movements(warehouse_id);
CREATE INDEX idx_stock_movements_created_at ON stock_movements(created_at);
CREATE INDEX idx_stock_movements_type ON stock_movements(movement_type);
