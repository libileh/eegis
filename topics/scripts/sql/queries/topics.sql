-- name: CreateTopic :one
INSERT INTO topics (id,
                    name,
                    description,
                    created_at)
VALUES ($1, $2, $3, $4) RETURNING
    id, created_at;

-- name: GetTopicByName :one
SELECT *
FROM topics
WHERE name = $1;

-- name: UpdateTopic :one
UPDATE topics
SET name        = $2,
    description = $3
WHERE id = $1 RETURNING
    id, name, description, created_at;

-- name: GetUserFollowedTopics :many
SELECT t.*
FROM topic_users tu
         JOIN topics t ON tu.topic_id = t.id
WHERE tu.user_id = $1;

-- name: GetAllTopics :many
SELECT *
FROM topics;

-- name: GetTopicById :one
SELECT *
FROM topics
WHERE id = $1;

-- name: FollowTopic :exec
INSERT INTO topic_users (topic_id, user_id, followed_at)
VALUES ($1, $2, $3);

-- name: UnfollowTopic :exec
DELETE FROM topic_users
WHERE topic_id = $1 AND user_id = $2;