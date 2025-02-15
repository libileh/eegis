package domain

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
)

/*
UserRepository defines the interface for user-related database operations.
This interface encapsulates all the methods required for managing user data
in the database, providing a clean separation between the database implementation
and business logic.
*/
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*uuid.UUID, *errors.CustomError)
	GetById(ctx context.Context, id uuid.UUID) (*User, *errors.CustomError)
	CreateAndInvite(ctx context.Context, user *User, userInvitation *UserInvitation) (*uuid.UUID, *errors.CustomError)
	ActivateUser(ctx context.Context, token string) *errors.CustomError
	DeleteUser(ctx context.Context, id uuid.UUID) *errors.CustomError
	GetByEmail(ctx context.Context, email string) (*User, *errors.CustomError)
}

type RoleRepository interface {
	GetByRoleName(ctx context.Context, role string) (*Role, *errors.CustomError)
	GetById(ctx context.Context, id int16) (*Role, *errors.CustomError)
}

type FollowerRepository interface {
	FollowUser(ctx context.Context, follower *Follower) (*Follower, *errors.CustomError)
	UnfollowUser(ctx context.Context, follower *Follower) *errors.CustomError
}
