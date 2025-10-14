CREATE TABLE IF NOT EXISTS soft_to_technic (
    soft_id BIGINT NOT NULL REFERENCES soft(id) ON DELETE CASCADE,
    technic_id BIGINT NOT NULL REFERENCES technic(id) ON DELETE CASCADE,
    PRIMARY KEY (soft_id, technic_id)
)