package domain

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
)

type TopicRepository interface {
	Create(ctx context.Context, topic *Topic) (*uuid.UUID, *errors.CustomError)
	UpdateTopic(ctx context.Context, id uuid.UUID, version int32, topic *Topic) (*Topic, *errors.CustomError)
	GetTopicByName(ctx context.Context, name string) (*Topic, *errors.CustomError)
	GetDBFollowerTopics(ctx context.Context, followerId uuid.UUID) ([]Topic, *errors.CustomError)
	GetAllTopics(ctx context.Context) ([]Topic, *errors.CustomError)
	FollowTopic(ctx context.Context, id uuid.UUID, id2 uuid.UUID) *errors.CustomError
}
