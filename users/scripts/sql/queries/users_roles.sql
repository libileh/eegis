-- name: GetByRoleName :one
SELECT *
FROM users_roles
WHERE name = $1;

-- name: GetById :one
SELECT *
FROM users_roles
WHERE id = $1;

