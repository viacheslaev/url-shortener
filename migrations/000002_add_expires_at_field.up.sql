ALTER TABLE links
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;

CREATE INDEX links_expires_idx
    ON links(expires_at)
    WHERE expires_at IS NOT NULL;