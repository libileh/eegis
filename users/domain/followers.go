package domain

import (
	"github.com/google/uuid"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"

	"time"
)

type Topic struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"title"`
	Description string    `json:"description"`
}

// swagger:model Follower
type Follower struct {
	UserId     uuid.UUID `json:"user_id"`
	FollowerId uuid.UUID `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func DbToFollower(dbFollower *database.Follower) *Follower {
	return &Follower{
		UserId:     dbFollower.UserID,
		FollowerId: dbFollower.FollowerID,
		CreatedAt:  dbFollower.CreatedAt,
	}
}
