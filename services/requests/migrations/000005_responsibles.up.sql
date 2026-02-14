CREATE TABLE IF NOT EXISTS responsibles
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL
);