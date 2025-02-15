-- +goose Up
CREATE EXTENSION IF NOT EXISTS "citext";
-- Enables case-insensitive text
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    password bytea NOT NULL, -- Allow flexibility for hashed passwords
    email CITEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Default to current time
);

-- +goose Down
DROP TABLE users IF EXISTS;