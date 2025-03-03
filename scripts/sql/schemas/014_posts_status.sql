-- +goose Up
ALTER TABLE posts
    ADD COLUMN status VARCHAR(10) CHECK (status IN ('pending', 'approved', 'rejected')) NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE posts DROP COLUMN status;
