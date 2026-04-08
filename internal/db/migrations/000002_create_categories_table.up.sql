-- Creates hierarchical categories table for product organization
-- parent_id = 0 means root category (top level)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_id INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast category tree queries
CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- Insert root category
INSERT INTO categories (id, name, parent_id, description) VALUES 
(1, 'Electronics', 0, 'All electronic products')
ON CONFLICT (id) DO NOTHING;
