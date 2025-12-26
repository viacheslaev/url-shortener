CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_links_created_at ON links (created_at);
