CREATE TABLE IF NOT EXISTS request_status_history
(
    id          SERIAL PRIMARY KEY,
    request_id  UUID                     NOT NULL REFERENCES requests (id) ON DELETE CASCADE,
    status      request_status           NOT NULL,
    comment     TEXT,
    changed_by  UUID REFERENCES users (id) ON DELETE SET NULL,
    changed_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_request_status_history_request_id ON request_status_history (request_id);

