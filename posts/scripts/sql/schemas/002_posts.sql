-- +goose Up
CREATE TABLE IF NOT EXISTS posts
(
    id         UUID PRIMARY KEY,
    title      TEXT      NOT NULL,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds IF EXISTS;