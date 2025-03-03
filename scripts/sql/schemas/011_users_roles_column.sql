-- +goose Up
ALTER TABLE
    IF EXISTS users
    ADD
    COLUMN role_id INT REFERENCES users_roles(id) DEFAULT 1;


-- +goose Down
ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS role_id;