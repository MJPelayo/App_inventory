-- Physical warehouse locations with capacity tracking
-- current_occupancy tracks total units currently stored
CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(200),
    capacity INTEGER DEFAULT 0 CHECK (capacity >= 0),
    current_occupancy INTEGER DEFAULT 0 CHECK (current_occupancy >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for warehouse location searches
CREATE INDEX idx_warehouses_location ON warehouses(location);