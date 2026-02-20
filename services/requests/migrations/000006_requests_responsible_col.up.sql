ALTER TABLE requests DROP COLUMN IF EXISTS responsible_info;
ALTER TABLE requests ADD COLUMN IF NOT EXISTS responsible_id uuid;
ALTER TABLE requests ADD FOREIGN KEY (responsible_id) REFERENCES responsibles(id);
CREATE INDEX ON requests(responsible_id);