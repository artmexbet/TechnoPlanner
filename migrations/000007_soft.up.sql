CREATE TABLE IF NOT EXISTS soft (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    install_link TEXT NOT NULL,
    description TEXT
);