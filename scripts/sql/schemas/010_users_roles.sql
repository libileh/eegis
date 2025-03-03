-- +goose UP
CREATE TABLE IF NOT EXISTS users_roles(
    id SMALLINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level SMALLINT NOT NULL DEFAULT 0,
    description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users_roles;
