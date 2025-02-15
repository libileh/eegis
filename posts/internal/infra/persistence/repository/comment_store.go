package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

type CommentStore struct {
	Queries *database.Queries
}

func (c *CommentStore) Create(ctx context.Context, comment *domain.Comment) (*uuid.UUID, *errors.CustomError) {
	// Implementation for creating a post
	id, err := c.Queries.CreateComment(ctx, database.CreateCommentParams{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
	})

	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	return &id, nil
}
