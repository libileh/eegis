-- +goose Up
CREATE TABLE IF NOT EXISTS users_invitations (
    token VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL,
    expiry TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE IF EXISTS users_invitations;
