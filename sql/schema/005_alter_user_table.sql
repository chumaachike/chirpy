-- +goose Up
ALTER TABLE users
    ADD CONSTRAINT email_not_empty CHECK (length(trim(email)) > 0);

ALTER TABLE users
    ADD CONSTRAINT password_not_empty CHECK (length(trim(hashed_password)) > 0);

-- +goose Down
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS email_not_empty;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS password_not_empty;
