package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/application/mapper"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"
)

type UserStore struct {
	Queries *database.Queries
	DB      *sql.DB
}

func (u *UserStore) GetById(ctx context.Context, id uuid.UUID) (*domain.User, *errors.CustomError) {
	dbUser, err := u.Queries.GetUserById(ctx, id)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return mapper.DbToUser(dbUser), nil
}

// executeInTransaction wraps a function within a database transaction.
func (u *UserStore) executeInTransaction(ctx context.Context, fn func(q *database.Queries) *errors.CustomError) *errors.CustomError {
	tx, err := u.DB.Begin()
	if err != nil {
		return errors.HandleDBError(err)
	}

	defer tx.Rollback()

	qtx := u.Queries.WithTx(tx)
	if customErr := fn(qtx); customErr != nil {
		return customErr
	}

	if err := tx.Commit(); err != nil {
		return errors.HandleDBError(err)
	}

	return nil
}

// createUser is a generic helper to handle user creation logic.
func (u *UserStore) createUser(ctx context.Context, q *database.Queries, user *domain.User) (*uuid.UUID, *errors.CustomError) {
	dbUser, err := q.CreateUser(ctx, database.CreateUserParams{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.PasswordHash,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		RoleID:    user.RoleID,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return &dbUser.ID, nil
}

// CreateUser creates a user without using a transaction.
func (u *UserStore) CreateUser(ctx context.Context, user *domain.User) (*uuid.UUID, *errors.CustomError) {
	return u.createUser(ctx, u.Queries, user)
}

// CreateAndInvite creates a user and sends an invitation, wrapped in a transaction.
func (u *UserStore) CreateAndInvite(
	ctx context.Context,
	user *domain.User,
	userInvitation *domain.UserInvitation) (*uuid.UUID, *errors.CustomError) {

	var userID *uuid.UUID
	err := u.executeInTransaction(ctx, func(qtx *database.Queries) *errors.CustomError {
		// Create user within transaction
		var customErr *errors.CustomError
		userID, customErr = u.createUser(ctx, qtx, user)
		if customErr != nil {
			return customErr
		}

		// Create user invitation within transaction
		_, err := qtx.CreateUserInvitation(ctx, database.CreateUserInvitationParams{
			Token:  userInvitation.Token,
			UserID: *userID,
			Expiry: userInvitation.Expiry,
		})
		if err != nil {
			return errors.HandleDBError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return userID, nil
}

func (u *UserStore) ActivateUser(ctx context.Context, token string) *errors.CustomError {

	return u.executeInTransaction(ctx, func(q *database.Queries) *errors.CustomError {
		//1. find user by token
		user, err := u.Queries.GetUserFromUserInvitation(ctx, token)
		if err != nil {
			return errors.HandleDBError(err)
		}
		//2. update user activate to true
		user.IsActive = true
		_, err = u.Queries.UpdateUser(ctx, database.UpdateUserParams{
			ID:       user.ID,
			IsActive: user.IsActive,
			//todo using a generic update isn't the best scenario
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		})

		if err != nil {
			return errors.HandleDBError(err)
		}

		//3. clean the invitation delete by token
		if err = u.Queries.DeleteUserInvitation(ctx, user.ID); err != nil {
			return errors.HandleDBError(err)
		}

		return nil
	})
}

// DeleteUser : SAGA Pattern for Rollback user CreateAndInvite process if Mailer.Send fails
func (u *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) *errors.CustomError {
	return u.executeInTransaction(ctx, func(q *database.Queries) *errors.CustomError {
		//1. Delete user invitation
		if err := u.Queries.DeleteUserInvitation(ctx, id); err != nil {
			return errors.HandleDBError(err)
		}

		//2. Delete user
		if err := u.Queries.DeleteUser(ctx, id); err != nil {
			return errors.HandleDBError(err)
		}
		return nil
	})
}

func (u *UserStore) GetByEmail(ctx context.Context, email string) (*domain.User, *errors.CustomError) {
	dbUserRow, err := u.Queries.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	user := domain.User{
		ID:        dbUserRow.ID,
		Username:  dbUserRow.Username,
		Email:     dbUserRow.Email,
		IsActive:  dbUserRow.IsActive,
		CreatedAt: dbUserRow.CreatedAt,
		RoleID:    dbUserRow.RoleID,
	}
	return &user, nil
}
