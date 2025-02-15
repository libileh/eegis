-- name: CreateUser :one
INSERT INTO users (id,
                   username,
                   password,
                   email,
                   created_at,
                   role_id)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
  AND is_active = true;

-- name: CreateUserInvitation :one
INSERT INTO users_invitations(token, user_id, expiry)
VALUES ($1, $2, $3) RETURNING *;


-- name: GetUserFromUserInvitation :one
SELECT u.*
FROM users u
         JOIN users_invitations ui
              ON u.id = ui.user_id
WHERE ui.token = $1;


-- name: UpdateUser :one
UPDATE users
SET username=$1,
    password=$2,
    email=$3,
    is_active=$4
WHERE id = $5
RETURNING *;


-- name: DeleteUserInvitation :exec
DELETE FROM users_invitations WHERE user_id=$1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;

-- name: GetByEmail :one
SELECT id, username, email, password, created_at, is_active, role_id
FROM users
WHERE email=$1 and is_active=true;


