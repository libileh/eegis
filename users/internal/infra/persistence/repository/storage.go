package repository

import (
	"database/sql"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"
)

/*
Storage is the main interface that defines the contract for accessing repositories.
It serves as a facade pattern implementation, providing a unified interface to access
different types of repositories (Users and Posts) while hiding the complexity
of repository instantiation and management.
*/
type Storage interface {
	Users() domain.UserRepository
	Roles() domain.RoleRepository
	Followers() domain.FollowerRepository
}

/*
RepositoryStorage implements the Storage interface specifically for PostgreSQL database.
It maintains:
- A database connection instance
- User repository implementation
- Post repository implementation

This struct serves as a concrete implementation of the Storage interface,
specifically tailored for PostgreSQL database operations.
*/
type RepositoryStorage struct {
	CustomErrors *errors.CustomError
	DB           *sql.DB
	Queries      *database.Queries
	UserRepo     domain.UserRepository
	RoleRepo     domain.RoleRepository
	FollowerRepo domain.FollowerRepository
}

/*
Users returns the UserRepository implementation.
This method provides access to the user repository instance,
allowing for user-related database operations while maintaining
encapsulation of the actual repository implementation.
*/
func (u *RepositoryStorage) Users() domain.UserRepository {
	return u.UserRepo
}

func (r *RepositoryStorage) Roles() domain.RoleRepository         { return r.RoleRepo }
func (f *RepositoryStorage) Followers() domain.FollowerRepository { return f.FollowerRepo }

/*
NewPostgresStorage creates and initializes a new PostgresStorage instance.
It takes a database connection as input and sets up both user and post
repositories with this connection.

Parameters:
  - db: A pointer to sql.DB representing an active PostgreSQL database connection

Returns:
  - A pointer to a fully configured PostgresStorage instance
*/

func NewPostgresStorage(db *sql.DB) *RepositoryStorage {

	queries := database.New(db)
	return &RepositoryStorage{
		DB:           db,
		Queries:      queries,
		UserRepo:     &UserStore{Queries: queries, DB: db},
		RoleRepo:     &RoleStore{Queries: queries},
		FollowerRepo: &FollowStore{Queries: queries},
	}
}
