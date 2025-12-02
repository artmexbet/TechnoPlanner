ALTER TABLE users
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

ALTER TABLE requests
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS updated_by UUID REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_requests_deleted_at ON requests (deleted_at);
CREATE INDEX IF NOT EXISTS idx_requests_created_by ON requests (created_by);
CREATE INDEX IF NOT EXISTS idx_requests_updated_by ON requests (updated_by);

ALTER TABLE equipment
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS updated_by UUID REFERENCES users (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_equipment_deleted_at ON equipment (deleted_at);
CREATE INDEX IF NOT EXISTS idx_equipment_created_by ON equipment (created_by);
CREATE INDEX IF NOT EXISTS idx_equipment_updated_by ON equipment (updated_by);

