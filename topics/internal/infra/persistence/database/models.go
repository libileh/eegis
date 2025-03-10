// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
}

type TopicUser struct {
	TopicID    uuid.UUID
	UserID     uuid.UUID
	FollowedAt time.Time
}
