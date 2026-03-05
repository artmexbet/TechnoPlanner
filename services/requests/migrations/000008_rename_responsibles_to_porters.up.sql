ALTER TABLE requests DROP CONSTRAINT IF EXISTS requests_responsible_id_fkey;
ALTER TABLE requests RENAME COLUMN responsible_id TO porter_id;
ALTER TABLE responsibles RENAME TO porters;
ALTER TABLE requests ADD CONSTRAINT requests_porter_id_fkey FOREIGN KEY (porter_id) REFERENCES porters(id);

