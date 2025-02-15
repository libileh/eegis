-- name: CreateComment :one
INSERT INTO
    comments (
        id,
        content,
        post_id,
        user_id,
        created_at
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    id;

-- name: GetCommentsByPostID :many
SELECT 
    c.id,
    c.post_id,
    c.user_id as commentUserId,
    c.content,
    c.created_at,
    u.username,
    u.id as userId
FROM comments c
INNER JOIN users u ON u.id = c.user_id
WHERE c.post_id = $1
ORDER BY c.created_at DESC;