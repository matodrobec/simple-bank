CREATE TABLE IF NOT EXISTS verify_emails (
    id bigserial PRIMARY KEY,
    username varchar NOT NULL,
    email varchar NOT NULL,
    secret_code varchar NOT NULL,
    is_used boolean NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_at TIMESTAMPTZ NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE verify_emails
ADD FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_email_verified bool NOT NULL DEFAULT false;