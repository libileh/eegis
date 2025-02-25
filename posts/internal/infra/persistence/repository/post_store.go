package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
	"log"
	"time"
)

/*** TODO replace domain by DTO*/

type PostStore struct {
	Queries *database.Queries
}

func (p *PostStore) Create(ctx context.Context, post *domain.Post) (*uuid.UUID, *errors.CustomError) {
	// Implementation for creating a post
	dbPost, err := p.Queries.CreatePost(ctx, database.CreatePostParams{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UserID:    post.UserID,
	})

	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	return &dbPost.ID, nil
}

func (p *PostStore) GetPostById(ctx context.Context, id uuid.UUID) (*domain.Post, *errors.CustomError) {
	// Retrieve the post from the database
	dbPost, err := p.Queries.GetPostById(ctx, id)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	// Convert the database post to a domain post
	post := domain.DbToPost(dbPost)

	// Retrieve and set the comments for the post
	comments, customErr := p.GetCommentsByPostId(ctx, id)
	if customErr != nil {
		return nil, customErr // Already handled by GetCommentsByPostId
	}
	post.Comments = comments

	return post, nil
}

func (p *PostStore) GetCommentsByPostId(ctx context.Context, postId uuid.UUID) ([]domain.Comment, *errors.CustomError) {
	// Retrieve comments from the database
	dbComments, err := p.Queries.GetCommentsByPostID(ctx, postId)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	// Convert database comments to domain comments
	return domain.DbToComments(dbComments), nil
}

func (p *PostStore) GetAllPosts(ctx context.Context) ([]domain.Post, *errors.CustomError) {

	posts, err := p.Queries.GetAllPosts(ctx)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return domain.DbToPosts(posts), nil
}

func (p *PostStore) DeletePost(ctx context.Context, id uuid.UUID) *errors.CustomError {

	if err := p.Queries.DeletePost(ctx, id); err != nil {
		return errors.HandleDBError(err)
	}
	return nil
}

func (p *PostStore) UpdatePost(ctx context.Context, id uuid.UUID, version int32, post *domain.Post) (*domain.UpdateResponse, *errors.CustomError) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel() // Ensure the cancel function is called to release resources
	select {
	case <-ctx.Done():
		log.Printf("Context canceled before query: %v", ctx.Err())
		return nil, errors.NewCustomError(errors.InternalServer, "context canceled before query execution")
	default:
	}

	res, err := p.Queries.UpdatePost(ctx, database.UpdatePostParams{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
		ID:      id,
		Version: version,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return &domain.UpdateResponse{
		ID:      res.ID,
		Version: res.Version,
	}, nil
}

func (p *PostStore) GetFeed(ctx context.Context, userId uuid.UUID, paginate domain.Paginated) ([]domain.Feed, *errors.CustomError) {
	res, err := p.Queries.GetUserFeed(ctx, database.GetUserFeedParams{
		UserID:  userId,
		Column2: domain.MapToNullString(&paginate.Search),
		Column3: domain.MapToNullString(&paginate.Search),
		Tags:    paginate.Tags,
		Column5: domain.SortPaginateToBool(paginate.Sort),
		Limit:   paginate.Limit,
		Offset:  paginate.Offset,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return domain.DbToUserFeed(res), nil
}
