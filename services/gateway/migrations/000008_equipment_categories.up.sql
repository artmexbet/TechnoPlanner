CREATE TABLE IF NOT EXISTS equipment_categories
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255)             NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE,
    created_by  UUID REFERENCES users (id) ON DELETE SET NULL,
    updated_by  UUID REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS equipment_category_links
(
    equipment_id INT NOT NULL REFERENCES equipment (id) ON DELETE CASCADE,
    category_id  INT NOT NULL REFERENCES equipment_categories (id) ON DELETE CASCADE,
    PRIMARY KEY (equipment_id, category_id)
);

