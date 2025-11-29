CREATE TYPE request_status AS ENUM (
    'canceled',
    'pending',
    'assigned',
    'in_progress',
    'completed',
    'rejected'
    );


CREATE TABLE IF NOT EXISTS requests
(
    id                  UUID PRIMARY KEY,
    telegram_user_info  jsonb,
    request_text        TEXT,
    status              request_status              NOT NULL DEFAULT 'pending',
    schedule_time       TEXT                        NOT NULL,
    end_time            TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    address             TEXT                        NOT NULL,
    responsible_user_id UUID                        REFERENCES users (id) ON DELETE SET NULL,
    created_at          TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_requests_status ON requests (status);
CREATE INDEX IF NOT EXISTS idx_requests_responsible_user_id ON requests (responsible_user_id);
