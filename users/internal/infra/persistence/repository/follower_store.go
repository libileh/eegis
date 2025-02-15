package repository

import (
	"context"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"
)

type FollowStore struct {
	Queries *database.Queries
}

func (f *FollowStore) FollowUser(ctx context.Context, follower *domain.Follower) (*domain.Follower, *errors.CustomError) {
	dbFollower, err := f.Queries.Follow(ctx, database.FollowParams{
		UserID:     follower.UserId,
		FollowerID: follower.FollowerId,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return domain.DbToFollower(&dbFollower), nil
}

func (f *FollowStore) UnfollowUser(ctx context.Context, follower *domain.Follower) *errors.CustomError {
	err := f.Queries.Unfollow(ctx, database.UnfollowParams{
		UserID:     follower.UserId,
		FollowerID: follower.FollowerId,
	})
	if err != nil {
		return errors.HandleDBError(err)
	}
	return nil
}
