-- Removes the reserved_quantity column and related constraints
DROP INDEX IF EXISTS idx_inventory_available_quantity;
ALTER TABLE inventory DROP CONSTRAINT IF EXISTS chk_reserved_quantity;
ALTER TABLE inventory DROP COLUMN IF EXISTS reserved_quantity;