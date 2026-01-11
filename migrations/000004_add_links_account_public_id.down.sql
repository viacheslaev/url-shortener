DROP INDEX IF EXISTS idx_links_account_public_id;

ALTER TABLE links
    DROP CONSTRAINT IF EXISTS links_account_public_id_fkey;

ALTER TABLE links
    DROP COLUMN IF EXISTS account_public_id;
