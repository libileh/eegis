// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id,
                   content,
                   title,
                   created_at,
                   user_id)
VALUES ($1, $2, $3, $4, $5) RETURNING
    id,
    created_at
`

type CreatePostParams struct {
	ID        uuid.UUID
	Content   string
	Title     string
	CreatedAt time.Time
	UserID    uuid.UUID
}

type CreatePostRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Content,
		arg.Title,
		arg.CreatedAt,
		arg.UserID,
	)
	var i CreatePostRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getAllPosts = `-- name: GetAllPosts :many
SELECT id, title, content, created_at, user_id, tags, version
FROM posts
`

func (q *Queries) GetAllPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UserID,
			pq.Array(&i.Tags),
			&i.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostById = `-- name: GetPostById :one
SELECT id, title, content, created_at, user_id, tags, version
FROM posts
WHERE id = $1
`

func (q *Queries) GetPostById(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UserID,
		pq.Array(&i.Tags),
		&i.Version,
	)
	return i, err
}

const getUserFeed = `-- name: GetFeed :many

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
    LIMIT $6 OFFSET $7
`

type GetUserFeedParams struct {
	UserID  uuid.UUID
	Column2 sql.NullString
	Column3 sql.NullString
	Tags    []string
	Column5 interface{}
	Limit   int32
	Offset  int32
}

type GetUserFeedRow struct {
	ID            uuid.UUID
	Title         string
	UserID        uuid.UUID
	Content       string
	Tags          []string
	Version       int32
	CreatedAt     time.Time
	Username      string
	TotalComments int64
}

// Parameters:
// $1: userId
// $2: title
// $3: content
// $4: tags
// $5: orderBy
// $6: limit
// $7: offset
func (q *Queries) GetUserFeed(ctx context.Context, arg GetUserFeedParams) ([]GetUserFeedRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeed,
		arg.UserID,
		arg.Column2,
		arg.Column3,
		pq.Array(arg.Tags),
		arg.Column5,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedRow
	for rows.Next() {
		var i GetUserFeedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.UserID,
			&i.Content,
			pq.Array(&i.Tags),
			&i.Version,
			&i.CreatedAt,
			&i.Username,
			&i.TotalComments,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET title   = $1,
    content = $2,
    tags    = $3,
    version = version + 1
WHERE id = $4
  and version = $5 RETURNING id, version
`

type UpdatePostParams struct {
	Title   string
	Content string
	Tags    []string
	ID      uuid.UUID
	Version int32
}

type UpdatePostRow struct {
	ID      uuid.UUID
	Version int32
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (UpdatePostRow, error) {
	row := q.db.QueryRowContext(ctx, updatePost,
		arg.Title,
		arg.Content,
		pq.Array(arg.Tags),
		arg.ID,
		arg.Version,
	)
	var i UpdatePostRow
	err := row.Scan(&i.ID, &i.Version)
	return i, err
}
