-- Tracks exact physical location of products within a warehouse
-- Uses aisle/side/shelf/layer system (no bins)
-- One product can be in multiple locations within the same warehouse
CREATE TABLE IF NOT EXISTS product_locations (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    aisle_number INTEGER NOT NULL CHECK (aisle_number > 0),
    side VARCHAR(5) NOT NULL CHECK (side IN ('left', 'right')),
    shelf_number INTEGER NOT NULL CHECK (shelf_number > 0),
    layer VARCHAR(10) NOT NULL CHECK (layer IN ('top', 'middle', 'middle2', 'middle3', 'bottom')),
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
    UNIQUE(product_id, warehouse_id, aisle_number, side, shelf_number, layer)
);

-- Indexes for location-based queries and picking optimization
CREATE INDEX idx_product_locations_warehouse ON product_locations(warehouse_id);
CREATE INDEX idx_product_locations_pick_path ON product_locations(warehouse_id, aisle_number, side, shelf_number);