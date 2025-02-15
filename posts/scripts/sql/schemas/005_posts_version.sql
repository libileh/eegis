-- +goose Up
ALTER TABLE posts ADD COLUMN version int NOT NULL default 0;

-- +goose Down
ALTER TABLE posts DROP COLUMN version;