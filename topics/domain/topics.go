package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/libileh/eegis/topics/internal/infra/request"
	"time"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Author    uuid.UUID `json:"author"`
	Tags      []string  `json:"tags"` //todo replace by topics
	Version   int32     `json:"version" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"name"`
	Email    string    `json:"email"`
}
type Topic struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"title"`
	Description string    `json:"description"`
	TopicPosts  []Post    `json:"topic_posts"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTopic(topicPayload request.TopicPayload) (*Topic, error) {
	if len(topicPayload.Name) == 0 || len(topicPayload.Description) == 0 {
		return nil, errors.New("both title and description are required")
	}
	if len(topicPayload.Name) > 100 {
		return nil, errors.New("title exceeds 100 character limit")
	}

	return &Topic{
		Id:          uuid.New(),
		Name:        topicPayload.Name,
		Description: topicPayload.Description,
		CreatedAt:   time.Now(),
	}, nil
}
