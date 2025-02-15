package domain

import (
	"github.com/google/uuid"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"

	"time"
)

// swagger:model Follower
type Follower struct {
	// The unique identifier of the user being followed
	UserId uuid.UUID `json:"user_id"`
	// The unique identifier of the follower
	FollowerId uuid.UUID `json:"follower_id"`
	// The time when the follow relationship was created
	CreatedAt time.Time `json:"created_at"`
}

func DbToFollower(dbFollower *database.Follower) *Follower {
	return &Follower{
		UserId:     dbFollower.UserID,
		FollowerId: dbFollower.FollowerID,
		CreatedAt:  dbFollower.CreatedAt,
	}
}
