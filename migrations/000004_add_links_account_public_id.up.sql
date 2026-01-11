ALTER TABLE links
    ADD COLUMN IF NOT EXISTS account_public_id UUID;

ALTER TABLE links
    ADD CONSTRAINT links_account_public_id_fkey
        FOREIGN KEY (account_public_id) REFERENCES accounts (public_id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_links_account_public_id
    ON links(account_public_id);
