package cache

import (
	"context"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/domain"
	"github.com/redis/go-redis/v9"
)

type CacheStorage interface {
	CacheUsers() CacheUserRepository
}

type CacheUserRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.User, *errors.CustomError)
	Set(ctx context.Context, user *domain.User) *errors.CustomError
}

type RedisStorage struct {
	cacheDB        *redis.Client
	CacheUsersRepo CacheUserRepository
}

func (uc *RedisStorage) Users() CacheUserRepository { return uc.CacheUsersRepo }

func NewRedisStorage(cacheDB *redis.Client) *RedisStorage {
	return &RedisStorage{
		cacheDB:        cacheDB,
		CacheUsersRepo: &CacheUserStore{cacheDB: cacheDB},
	}
}
