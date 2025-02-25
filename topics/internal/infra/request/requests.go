package request

import (
	"github.com/google/uuid"
	"time"
)

type TopicPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type TopicDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}

type UserDTO struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	// Add only the fields needed for the Feed
}
