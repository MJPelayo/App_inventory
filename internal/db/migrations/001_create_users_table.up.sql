-- Creates the users table for storing all system user accounts
-- Supports four roles: admin, sales_manager, warehouse_manager, supply_manager
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'sales_manager', 'warehouse_manager', 'supply_manager')),
    department VARCHAR(100),
    sales_target DECIMAL(12,2),
    commission_rate DECIMAL(5,2) DEFAULT 5.0,
    warehouse_id INTEGER,
    shift VARCHAR(20),
    purchase_budget DECIMAL(12,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Index for faster email lookups during login
CREATE INDEX idx_users_email ON users(email);