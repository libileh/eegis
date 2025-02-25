package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/topics/domain"
	"github.com/libileh/eegis/topics/internal/client"
	"github.com/libileh/eegis/topics/internal/infra/request"
)

type Service struct {
	domain.TopicRepository
	*client.UserClient
}

func NewService(topicRepos domain.TopicRepository, userClient *client.UserClient) *Service {
	return &Service{
		TopicRepository: topicRepos,
		UserClient:      userClient,
	}
}

type UserService interface {
	GetUserById(id uuid.UUID) (*request.UserDTO, error)
}

// GetFollowerTopic retrieves the topics a user is following and wraps errors appropriately
// so that the API returns the correct HTTP status codes based on the error type.
func (s *Service) GetFollowerTopic(ctx context.Context, followerId uuid.UUID) ([]domain.Topic, *errors.CustomError) {
	// Retrieve topics followed by the user from the repository.
	topics, repoErr := s.TopicRepository.GetDBFollowerTopics(ctx, followerId)
	if repoErr != nil {
		return nil, repoErr
	}
	return topics, nil
}
