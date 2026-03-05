ALTER TABLE requests DROP CONSTRAINT IF EXISTS requests_porter_id_fkey;
ALTER TABLE requests RENAME COLUMN porter_id TO responsible_id;
ALTER TABLE porters RENAME TO responsibles;
ALTER TABLE requests ADD CONSTRAINT requests_responsible_id_fkey FOREIGN KEY (responsible_id) REFERENCES responsibles(id);

