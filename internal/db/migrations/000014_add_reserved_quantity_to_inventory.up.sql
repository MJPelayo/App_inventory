-- Adds reserved_quantity column to inventory table for order reservation system
-- This allows products to be "soft reserved" when an order is created but not yet fulfilled
-- Prevents overselling when multiple orders come in simultaneously
ALTER TABLE inventory 
ADD COLUMN reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0);

-- Add check constraint to ensure reserved_quantity never exceeds quantity
ALTER TABLE inventory 
ADD CONSTRAINT chk_reserved_quantity CHECK (reserved_quantity <= quantity);

-- Create index for faster availability queries
CREATE INDEX idx_inventory_available_quantity ON inventory(quantity - reserved_quantity);

-- Update existing records (reserved_quantity already defaults to 0, no data change needed)
COMMENT ON COLUMN inventory.reserved_quantity IS 'Quantity reserved for pending orders but not yet picked/shipped';