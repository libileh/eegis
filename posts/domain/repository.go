package domain

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
)

/*
PostRepository defines the interface for post-related database operations.
Similar to UserRepository, this interface provides a contract for managing
post data, ensuring consistent post management operations across different
database implementations.
*/
type PostRepository interface {
	Create(ctx context.Context, post *Post) (*uuid.UUID, *errors.CustomError)
	GetPostById(ctx context.Context, id uuid.UUID) (*Post, *errors.CustomError)
	GetAllPosts(ctx context.Context) ([]Post, *errors.CustomError)
	DeletePost(ctx context.Context, id uuid.UUID) *errors.CustomError
	UpdatePost(ctx context.Context, id uuid.UUID, version int32, post *Post) (*UpdateResponse, *errors.CustomError)
	GetCommentsByPostId(ctx context.Context, postId uuid.UUID) ([]Comment, *errors.CustomError)
}

/*
CommentRepository Same for Comments
*/
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) (*uuid.UUID, *errors.CustomError)
}

type FeedRepository interface {
	GetFeed(ctx context.Context, id uuid.UUID, paginate Paginated) ([]Feed, *errors.CustomError)
}
