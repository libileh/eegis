package client

import "github.com/google/uuid"

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	// Add only the fields needed for the Feed
}
