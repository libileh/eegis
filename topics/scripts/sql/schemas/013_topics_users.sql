-- +goose Up
-- Create junction table for the many-to-many relationship between topics and users
CREATE TABLE topic_users
(
    topic_id   UUID      NOT NULL,
    user_id    UUID      NOT NULL,
    followed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (topic_id, user_id),
    FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
-- Drop junction table first because of foreign key constraints
DROP TABLE IF EXISTS topic_users;