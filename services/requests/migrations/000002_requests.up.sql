CREATE TYPE request_status AS ENUM (
    'canceled',
    'pending',
    'assigned',
    'in_progress',
    'completed',
    'rejected'
);

CREATE TABLE IF NOT EXISTS requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_user_id UUID NOT NULL REFERENCES telegram_users(id) ON DELETE RESTRICT,
    request_text TEXT,
    status request_status NOT NULL DEFAULT 'pending',
    schedule_time TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_requests_telegram_user_id ON requests(telegram_user_id);
