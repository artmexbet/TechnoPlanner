ALTER TABLE equipment
    DROP CONSTRAINT IF EXISTS chk_reserved_le_quantity,
    DROP CONSTRAINT IF EXISTS chk_reserved_non_negative,
    DROP COLUMN IF EXISTS reserved_quantity,
    DROP COLUMN IF EXISTS quantity;

