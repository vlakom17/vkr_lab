CREATE TABLE charts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    title TEXT NOT NULL,
    genre TEXT NOT NULL,
    description TEXT,
    position_count INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);