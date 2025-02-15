package repository

import (
	"database/sql"
	"github.com/libileh/eegis/posts/domain"
	"github.com/libileh/eegis/posts/internal/infra/persistence/database"
)

/*
Storage is the main interface that defines the contract for accessing repositories.
It serves as a facade pattern implementation, providing a unified interface to access
different types of repositories (Users and Posts) while hiding the complexity
of repository instantiation and management.
*/
type Storage interface {
	Posts() domain.PostRepository
	Comments() domain.CommentRepository
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
	DB          *sql.DB
	Queries     *database.Queries
	PostRepo    domain.PostRepository
	CommentRepo domain.CommentRepository
	FeedRepo    domain.FeedRepository
}

/*
Posts returns the PostRepository implementation.
Similar to Users(), this method provides access to the post
repository instance while maintaining proper encapsulation
of the underlying implementation.
*/
func (p *RepositoryStorage) Posts() domain.PostRepository {
	return p.PostRepo
}

// Comments Same for Comments*/
func (c *RepositoryStorage) Comments() domain.CommentRepository {
	return c.CommentRepo
}

func (f *RepositoryStorage) Feeds() domain.FeedRepository {
	return f.FeedRepo
}

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
		DB:          db,
		Queries:     queries,
		PostRepo:    &PostStore{Queries: queries},
		CommentRepo: &CommentStore{Queries: queries},
	}
}
