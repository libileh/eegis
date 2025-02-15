-- +goose Up
ALTER TABLE posts ADD COLUMN tags VARCHAR(100) [];

-- +goose Down
ALTER TABLE posts DROP COLUMN tags;