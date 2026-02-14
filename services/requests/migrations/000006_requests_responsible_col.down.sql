ALTER TABLE requests DROP COLUMN IF EXISTS responsible_id;
ALTER TABLE requests ADD COLUMN IF NOT EXISTS responsible_info JSONB;