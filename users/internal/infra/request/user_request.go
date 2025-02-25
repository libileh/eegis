package request

import (
	"github.com/google/uuid"
	"github.com/libileh/eegis/users/domain"
	"time"
)

// swagger:model UserRequest
type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	RoleID   int16  `json:"role_id" validate:"required,min=1,max=3"`
}

func MapToUser(payload UserRequest) domain.User {
	user := domain.User{
		ID:        uuid.New(),
		Username:  payload.Username,
		Email:     payload.Email,
		CreatedAt: time.Now(),
		RoleID:    payload.RoleID,
	}
	return user
}

type TopicDTO struct {
	ID          string `json:"topic_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
