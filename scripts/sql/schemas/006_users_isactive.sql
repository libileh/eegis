-- +goose Up
ALTER TABLE users
    ADD COLUMN is_active BOOL NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_active;