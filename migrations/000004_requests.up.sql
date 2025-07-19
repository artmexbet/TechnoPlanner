CREATE TABLE IF NOT EXISTS requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_r VARCHAR(255) NOT NULL,
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    usage_interval INTERVAL NOT NULL,
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    responsible_user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    days TIMESTAMPTZ[] NOT NULL,
)