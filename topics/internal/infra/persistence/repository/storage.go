package repository

import (
	"database/sql"
	"github.com/libileh/eegis/topics/domain"
	"github.com/libileh/eegis/topics/internal/infra/persistence/database"
)

type Storage interface {
	Topics() domain.TopicRepository
}

func (t *RepositoryStorage) Topics() domain.TopicRepository { return t.TopicRepo }

type RepositoryStorage struct {
	DB        *sql.DB
	Queries   *database.Queries
	TopicRepo domain.TopicRepository
}

func NewPostgresStorage(db *sql.DB) *RepositoryStorage {
	queries := database.New(db)
	return &RepositoryStorage{
		DB:        db,
		Queries:   queries,
		TopicRepo: &TopicStore{Queries: queries},
	}
}
