DROP INDEX IF EXISTS idx_equipment_updated_by;
DROP INDEX IF EXISTS idx_equipment_created_by;
DROP INDEX IF EXISTS idx_equipment_deleted_at;
ALTER TABLE equipment
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_requests_updated_by;
DROP INDEX IF EXISTS idx_requests_created_by;
DROP INDEX IF EXISTS idx_requests_deleted_at;
ALTER TABLE requests
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users
    DROP COLUMN IF EXISTS deleted_at;

