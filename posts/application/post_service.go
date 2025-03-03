package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/eventbus"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/pkg/client"
	"time"
)

type ServiceManager struct {
	PostService *PostService
	UserService *client.HttpUserService
}

func NewServiceManager(postService *PostService, userService *client.HttpUserService) *ServiceManager {
	return &ServiceManager{
		PostService: postService,
		UserService: userService,
	}
}

type PostService struct {
	PostRepo    domain.PostRepository
	CommentRepo domain.CommentRepository
	FeedRepo    domain.FeedRepository
	EventBus    *eventbus.EventBus
}

func NewPostService(postRepo domain.PostRepository, commentRepo domain.CommentRepository, feedRepo domain.FeedRepository, eventBus *eventbus.EventBus) *PostService {
	return &PostService{
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
		FeedRepo:    feedRepo,
		EventBus:    eventBus,
	}
}

// ReviewPost delegates the review operation to the domain aggregate,
// ensuring that business logic is encapsulated within the domain.
func (s *PostService) ReviewPost(ctx context.Context, postId *uuid.UUID, ctxUser *auth.CtxUser, status string) (*domain.PostReview, *errors.CustomError) {
	// Authorization: Only ADMIN or Moderator can review posts.
	if ctxUser.ContextRole != auth.ADMIN && ctxUser.ContextRole != auth.Moderator {
		return nil, errors.NewBadRequest("user is not authorized to perform this operation")
	}

	// Retrieve the post from the repository.
	post, customErr := s.PostRepo.GetPostById(ctx, *postId)
	if customErr != nil {
		return nil, customErr
	}

	// Delegate the review logic to the domain aggregate.
	reviewPost, err := post.Review(ctxUser, status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	// Persist the updated post (with new status and version) to the repository.
	_, customErr = s.PostRepo.UpdatePost(ctx, *postId, post.Version, post)
	if customErr != nil {
		return nil, customErr
	}

	event := eventbus.Event{
		Type:       "review-post",
		OccurredAt: time.Now(),
		Data:       domain.NewPostStatusChangedEvent(post.ID, post.Status, post.UserID),
	}

	if err := s.EventBus.Publish(event); err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	return reviewPost, nil
}

type UserService interface {
	CheckRolePrecedence(userID uuid.UUID, roleName string) (bool, error)
	GetUsername(ctx context.Context, userID uuid.UUID) (string, error)
}
