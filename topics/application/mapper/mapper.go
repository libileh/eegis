package mapper

import (
	"github.com/libileh/eegis/topics/domain"
	"github.com/libileh/eegis/topics/internal/infra/persistence/database"
	"github.com/libileh/eegis/topics/internal/infra/request"
)

func MapTopic(dbTopic database.Topic) *domain.Topic {
	return &domain.Topic{
		Id:          dbTopic.ID,
		Name:        dbTopic.Name,
		Description: dbTopic.Description,
		CreatedAt:   dbTopic.CreatedAt,
	}
}

func MapToUser(userDTO *request.UserDTO) *domain.User {
	return &domain.User{
		Id:       userDTO.Id,
		Username: userDTO.Username,
		Email:    userDTO.Email,
	}
}

func MapTopics(dbTopics []database.Topic) []domain.Topic {
	var topics []domain.Topic
	for _, dbTopic := range dbTopics {
		topics = append(topics, *MapTopic(dbTopic))
	}
	return topics
}
