-- +goose Up
-- Create topics table
CREATE TABLE topics
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP -- Default to current time
);

-- +goose Down
-- Drop topics table
DROP TABLE IF EXISTS topics;