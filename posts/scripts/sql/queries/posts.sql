-- name: CreatePost :one
INSERT INTO posts (id,
                   content,
                   title,
                   created_at,
                   user_id)
VALUES ($1, $2, $3, $4, $5) RETURNING
    id,
    created_at;

-- name: GetPostById :one
SELECT *
FROM posts
WHERE id = $1;

-- name: GetAllPosts :many
SELECT *
FROM posts;

-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts
SET title   = $1,
    content = $2,
    tags    = $3,
    version = version + 1
WHERE id = $4
  and version = $5 RETURNING id, version;

-- name: GetUserFeed :many
-- Parameters:
-- $1: userId
-- $2: title
-- $3: content
-- $4: tags
-- $5: orderBy
-- $6: limit
-- $7: offset

SELECT
    p.id,
    p.title,
    p.user_id,
    p.content,
    p.tags,
    p.version,
    p.created_at,
    u.username,
    COUNT(DISTINCT c.id) AS total_comments
FROM posts p
         LEFT JOIN comments c ON c.post_id = p.id
         INNER JOIN users u ON p.user_id = u.id
         LEFT JOIN followers f ON f.follower_id = p.user_id AND f.user_id = $1
WHERE (f.follower_id IS NOT NULL OR p.user_id = $1)
  AND (
    p.title ILIKE '%' || $2 || '%'
        OR p.content ILIKE '%' || $3 || '%'
    )
  AND (
    p.tags @> $4 OR $4 = '{}'
    )
GROUP BY p.id, u.username
ORDER BY
    CASE WHEN $5 THEN p.created_at END DESC,
    CASE WHEN NOT $5 THEN p.created_at END ASC
    LIMIT $6 OFFSET $7;


