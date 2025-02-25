package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/topics/application/mapper"
	"github.com/libileh/eegis/topics/domain"
	"github.com/libileh/eegis/topics/internal/infra/persistence/database"
	"time"
)

type TopicStore struct {
	Queries *database.Queries
}

func (t *TopicStore) Create(ctx context.Context, topic *domain.Topic) (*uuid.UUID, *errors.CustomError) {
	dbTopic, err := t.Queries.CreateTopic(ctx, database.CreateTopicParams{
		ID:          topic.Id,
		Name:        topic.Name,
		Description: topic.Description,
		CreatedAt:   topic.CreatedAt,
	})
	if err != nil {
		return nil, errors.HandleDBError(err)
	}

	return &dbTopic.ID, nil
}

func (t *TopicStore) GetTopicByName(ctx context.Context, name string) (*domain.Topic, *errors.CustomError) {
	dbTopic, err := t.Queries.GetTopicByName(ctx, name)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	topic := mapper.MapTopic(dbTopic)
	return topic, nil
}

func (t *TopicStore) GetDBFollowerTopics(ctx context.Context, followerId uuid.UUID) ([]domain.Topic, *errors.CustomError) {
	dbTopics, err := t.Queries.GetUserFollowedTopics(ctx, followerId)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	topics := mapper.MapTopics(dbTopics)
	return topics, nil
}

func (t *TopicStore) GetAllTopics(ctx context.Context) ([]domain.Topic, *errors.CustomError) {
	dbTopics, err := t.Queries.GetAllTopics(ctx)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	topics := mapper.MapTopics(dbTopics)
	return topics, nil
}

func (t *TopicStore) FollowTopic(ctx context.Context, userId uuid.UUID, topicId uuid.UUID) *errors.CustomError {
	if err := t.Queries.FollowTopic(ctx, database.FollowTopicParams{
		TopicID:    topicId,
		UserID:     userId,
		FollowedAt: time.Now(),
	}); err != nil {
		return errors.HandleDBError(err)
	}
	return nil
}
func (t *TopicStore) UpdateTopic(ctx context.Context, id uuid.UUID, version int32, topic *domain.Topic) (*domain.Topic, *errors.CustomError) {
	return nil, nil
}
