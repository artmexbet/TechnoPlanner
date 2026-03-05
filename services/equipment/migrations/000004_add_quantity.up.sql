ALTER TABLE equipment
    ADD COLUMN IF NOT EXISTS quantity          INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reserved_quantity INT NOT NULL DEFAULT 0;

ALTER TABLE equipment
    ADD CONSTRAINT chk_reserved_non_negative CHECK (reserved_quantity >= 0),
    ADD CONSTRAINT chk_reserved_le_quantity  CHECK (reserved_quantity <= quantity);

