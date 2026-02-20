CREATE TABLE IF NOT EXISTS raw_requests
(
    id                   UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    telegram_id          BIGINT      NOT NULL,
    username             TEXT        NOT NULL DEFAULT '',
    first_name           TEXT        NOT NULL DEFAULT '',
    last_name            TEXT,
    raw_text             TEXT        NOT NULL,
    status               TEXT        NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'processed')),
    processed_request_id UUID REFERENCES requests (id) ON DELETE SET NULL,
    created_at           TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_raw_requests_telegram_id ON raw_requests (telegram_id);
CREATE INDEX IF NOT EXISTS idx_raw_requests_status ON raw_requests (status);

