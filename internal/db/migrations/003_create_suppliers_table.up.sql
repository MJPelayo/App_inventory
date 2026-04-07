-- Stores supplier information including contact details and performance metrics
-- Rating is calculated based on on-time deliveries (0-5 scale)
CREATE TABLE IF NOT EXISTS suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    contact_person VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    tax_id VARCHAR(50),
    payment_terms VARCHAR(50) DEFAULT 'Net 30',
    lead_time_days INTEGER DEFAULT 7,
    minimum_order INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    on_time_deliveries INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for supplier name searches
CREATE INDEX idx_suppliers_name ON suppliers(name);