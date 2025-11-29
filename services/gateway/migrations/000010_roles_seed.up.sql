INSERT INTO roles (id, name, description)
VALUES
    (1, 'admin', 'Администратор'),
    (2, 'porter', 'Носильщик')
ON CONFLICT (id) DO NOTHING;

