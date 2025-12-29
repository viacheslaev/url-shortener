CREATE TABLE accounts
(
    id            BIGSERIAL PRIMARY KEY,
    public_id     UUID         NOT NULL DEFAULT gen_random_uuid(),

    email         VARCHAR(320) NOT NULL,
    password_hash TEXT         NOT NULL,

    is_active     BOOLEAN      NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT accounts_public_id_unique UNIQUE (public_id)
);

CREATE UNIQUE INDEX accounts_email_active_idx
    ON accounts (email) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX accounts_public_id_idx
    ON accounts (public_id) WHERE deleted_at IS NULL