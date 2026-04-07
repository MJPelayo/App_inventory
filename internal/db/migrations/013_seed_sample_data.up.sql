-- ============================================
-- SAMPLE DATA FOR TEST 1 (2+ rows per table)
-- ============================================

-- Insert 4 sample users (one for each role)
INSERT INTO users (id, name, email, password_hash, role, department, sales_target, commission_rate, purchase_budget) VALUES
(1, 'John Admin', 'admin@system.com', 'admin123_hash', 'admin', 'IT', NULL, NULL, NULL),
(2, 'Sarah Sales', 'sarah@company.com', 'sales123_hash', 'sales_manager', 'Sales', 100000, 5.0, NULL),
(3, 'Mike Warehouse', 'mike@company.com', 'warehouse123_hash', 'warehouse_manager', 'Warehouse', NULL, NULL, NULL),
(4, 'Lisa Supply', 'lisa@company.com', 'supply123_hash', 'supply_manager', 'Supply', NULL, NULL, 500000);

-- Insert hierarchical categories (Electronics > Phones > Samsung/Apple)
INSERT INTO categories (id, name, parent_id, description) VALUES
(1, 'Electronics', 0, 'All electronic products'),
(2, 'Phones', 1, 'Mobile phones and smartphones'),
(3, 'Samsung', 2, 'Samsung brand phones'),
(4, 'Apple', 2, 'Apple brand phones'),
(5, 'Accessories', 1, 'Phone accessories and chargers'),
(6, 'Hardware', 0, 'Tools and hardware supplies'),
(7, 'Hammers', 6, 'All types of hammers'),
(8, 'Rubber Hammers', 7, 'Non-marring rubber mallets'),
(9, 'Claw Hammers', 7, 'Traditional claw hammers'),
(10, 'Plumbing', 6, 'Plumbing supplies');

-- Insert 4 sample suppliers
INSERT INTO suppliers (id, name, contact_person, phone, email, rating, total_orders, on_time_deliveries) VALUES
(1, 'Samsung Electronics', 'Kim Lee', '555-0101', 'samsung@supplier.com', 4.8, 10, 9),
(2, 'Apple Inc', 'Tim Cook', '555-0202', 'apple@supplier.com', 4.9, 15, 14),
(3, 'Stanley Tools', 'John Stanley', '555-0303', 'stanley@tools.com', 4.7, 8, 8),
(4, 'DeWalt Industrial', 'Tom DeWalt', '555-0404', 'dewalt@tools.com', 4.9, 12, 11);

-- Insert 8 sample products (Samsung, Apple, and Hardware)
INSERT INTO products (id, name, description, sku, brand, price, cost, category_id, supplier_id) VALUES
(101, 'Galaxy S23', 'Latest Samsung flagship phone with 256GB storage', 'SAM-GS23', 'Samsung', 999.99, 750.00, 3, 1),
(102, 'Galaxy Tab S9', 'Samsung premium tablet with S Pen included', 'SAM-TABS9', 'Samsung', 799.99, 600.00, 3, 1),
(103, 'iPhone 15', 'Latest Apple iPhone with dynamic island', 'APL-IP15', 'Apple', 1099.99, 800.00, 4, 2),
(104, 'MacBook Pro', 'Apple professional laptop with M3 chip', 'APL-MBP', 'Apple', 1999.99, 1500.00, 4, 2),
(105, 'Rubber Mallet', 'Non-marring rubber hammer for delicate work', 'STN-RM1', 'Stanley', 24.99, 15.00, 8, 3),
(106, 'Claw Hammer 16oz', 'Traditional claw hammer with fiberglass handle', 'STN-CH16', 'Stanley', 19.99, 12.00, 9, 3),
(107, 'Cordless Drill', '20V Max cordless drill with battery', 'DWT-CD20', 'DeWalt', 149.99, 110.00, 6, 4),
(108, 'Pipe Wrench', 'Heavy-duty pipe wrench for plumbing', 'DWT-PW18', 'DeWalt', 39.99, 28.00, 10, 4);

-- Insert 3 sample warehouses
INSERT INTO warehouses (id, name, location, capacity, current_occupancy) VALUES
(101, 'Main Warehouse', 'New York, NY', 10000, 53),
(102, 'West Coast Hub', 'Los Angeles, CA', 8000, 10),
(103, 'Central Distribution', 'Chicago, IL', 5000, 0);

-- Insert inventory records (product-warehouse combinations with quantities)
INSERT INTO inventory (id, product_id, warehouse_id, quantity, reorder_point, max_stock) VALUES
(1001, 101, 101, 25, 10, 100),
(1002, 102, 101, 15, 5, 50),
(1003, 103, 101, 8, 8, 50),    -- LOW STOCK (at reorder point)
(1004, 104, 101, 3, 3, 20),     -- LOW STOCK (at reorder point)
(1005, 101, 102, 10, 5, 50),
(1006, 105, 101, 50, 10, 200),
(1007, 106, 101, 35, 10, 150),
(1008, 107, 101, 12, 5, 60),
(1009, 108, 101, 8, 5, 40);

-- Insert product locations (aisle/side/shelf/layer tracking)
INSERT INTO product_locations (product_id, warehouse_id, aisle_number, side, shelf_number, layer, quantity) VALUES
(101, 101, 2, 'left', 1, 'top', 15),
(101, 101, 2, 'right', 3, 'middle', 10),
(102, 101, 2, 'left', 1, 'middle', 15),
(103, 101, 1, 'right', 2, 'middle', 8),
(104, 101, 1, 'right', 2, 'bottom', 3),
(105, 101, 4, 'left', 1, 'top', 30),
(105, 101, 4, 'left', 1, 'middle', 20),
(106, 101, 4, 'left', 2, 'top', 35),
(107, 101, 5, 'right', 1, 'top', 12),
(108, 101, 5, 'right', 1, 'bottom', 8);

-- Insert sample sales order
INSERT INTO sales_orders (id, order_number, customer_name, customer_email, customer_phone, shipping_address, delivery_type, status, subtotal, total_amount, created_by) VALUES
(301, 'SO-2024-001', 'John Smith', 'john.smith@email.com', '555-123-4567', '123 Main St, New York, NY 10001', 'delivery', 'pending', 3099.97, 3099.97, 2);

-- Insert order items for sales order
INSERT INTO order_items (id, order_id, order_type, product_id, product_name, quantity, unit_price) VALUES
(401, 301, 'sales', 101, 'Galaxy S23', 2, 999.99),
(402, 301, 'sales', 103, 'iPhone 15', 1, 1099.99);

-- Insert sample supply order (purchase order)
INSERT INTO supply_orders (id, po_number, supplier_id, status, order_date, expected_delivery_date, subtotal, total_amount, created_by) VALUES
(501, 'PO-2024-001', 1, 'processing', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '7 days', 15000.00, 15150.00, 4);

-- Insert stock movements (audit trail)
INSERT INTO stock_movements (product_id, warehouse_id, quantity_change, movement_type, reason, reference_number, performed_by, previous_quantity, new_quantity) VALUES
(101, 101, 25, 'received', 'Initial stock setup', 'INIT-001', 1, 0, 25),
(103, 101, 8, 'received', 'Initial stock setup', 'INIT-001', 1, 0, 8),
(105, 101, 50, 'received', 'Initial stock setup', 'INIT-001', 1, 0, 50),
(106, 101, 35, 'received', 'Initial stock setup', 'INIT-001', 1, 0, 35);

-- Insert audit log entries
INSERT INTO audit_log (user_id, action, entity_type, entity_id, ip_address) VALUES
(1, 'CREATE', 'user', 2, '127.0.0.1'),
(1, 'CREATE', 'user', 3, '127.0.0.1'),
(1, 'CREATE', 'user', 4, '127.0.0.1'),
(2, 'CREATE', 'sales_order', 301, '127.0.0.1');
