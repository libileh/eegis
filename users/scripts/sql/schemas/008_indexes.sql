-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- +goose Down
DROP INDEX IF EXISTS idx_users_username;
