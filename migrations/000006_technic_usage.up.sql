CREATE TABLE IF NOT EXISTS technic_usage (
    request_id UUID NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    technic_id BIGINT NOT NULL REFERENCES technic(id) ON DELETE CASCADE,
    PRIMARY KEY (request_id, technic_id)
);