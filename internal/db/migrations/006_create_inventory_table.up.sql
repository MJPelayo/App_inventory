-- Tracks stock levels per product per warehouse (separate from product catalog)
-- This is the SINGLE SOURCE OF TRUTH for all stock levels
CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reorder_point INTEGER DEFAULT 0 CHECK (reorder_point >= 0),
    max_stock INTEGER DEFAULT 0 CHECK (max_stock >= 0),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
    UNIQUE(product_id, warehouse_id)  -- One record per product per warehouse
);

-- Indexes for fast inventory queries
CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse_id ON inventory(warehouse_id);
CREATE INDEX idx_inventory_low_stock ON inventory(quantity, reorder_point);