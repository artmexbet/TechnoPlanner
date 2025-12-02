CREATE TABLE IF NOT EXISTS equipment_to_requests
(
    request_id   UUID                     NOT NULL REFERENCES requests (id) ON DELETE CASCADE,
    equipment_id INT                      NOT NULL REFERENCES equipment (id) ON DELETE CASCADE,
    quantity     INT                      NOT NULL DEFAULT 1,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (request_id, equipment_id)
);