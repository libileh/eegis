package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/pkg/client"
)

// Implementation of CombinedService
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
}

func NewPostService(postRepo domain.PostRepository, commentRepo domain.CommentRepository, feedRepo domain.FeedRepository) *PostService {
	return &PostService{
		PostRepo:    postRepo,
		CommentRepo: commentRepo,
		FeedRepo:    feedRepo,
	}
}

type UserService interface {
	CheckRolePrecedence(userID uuid.UUID, roleName string) (bool, error)
	GetUsername(ctx context.Context, userID uuid.UUID) (string, error)
}
