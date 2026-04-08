-- Remove all sample data (order matters due to foreign keys)
DELETE FROM audit_log;
DELETE FROM stock_movements;
DELETE FROM order_items;
DELETE FROM sales_orders;
DELETE FROM supply_orders;
DELETE FROM product_locations;
DELETE FROM inventory;
DELETE FROM products;
DELETE FROM suppliers;
DELETE FROM categories;
DELETE FROM users;
DELETE FROM warehouses;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE categories_id_seq RESTART WITH 1;
ALTER SEQUENCE suppliers_id_seq RESTART WITH 1;
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE warehouses_id_seq RESTART WITH 1;
ALTER SEQUENCE inventory_id_seq RESTART WITH 1;
ALTER SEQUENCE product_locations_id_seq RESTART WITH 1;
ALTER SEQUENCE sales_orders_id_seq RESTART WITH 1;
ALTER SEQUENCE supply_orders_id_seq RESTART WITH 1;
ALTER SEQUENCE order_items_id_seq RESTART WITH 1;
ALTER SEQUENCE stock_movements_id_seq RESTART WITH 1;
ALTER SEQUENCE audit_log_id_seq RESTART WITH 1;