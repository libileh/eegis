package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	customErr "github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheUserStore struct {
	cacheDB *redis.Client
}

func (uc *CacheUserStore) Get(ctx context.Context, id uuid.UUID) (*domain.User, *customErr.CustomError) {
	cacheKey := fmt.Sprintf("user:%v", id)
	result, err := uc.cacheDB.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) { //err == redis.Nil: no cache value present
		return nil, nil
	}
	if err != nil {
		return nil, customErr.HandleDBError(err)
	}
	var user domain.User
	if result != "" {
		err = json.Unmarshal([]byte(result), &user)
		if err != nil {
			return nil, customErr.HandleDBError(err)
		}
	}
	return &user, nil
}

func (uc *CacheUserStore) Set(ctx context.Context, user *domain.User) *customErr.CustomError {
	cacheKey := fmt.Sprintf("user:%v", user.ID)
	result, err := json.Marshal(user)
	if err != nil {
		return customErr.HandleDBError(err)
	}
	err = uc.cacheDB.SetEx(ctx, cacheKey, string(result), time.Minute).Err()
	if err != nil {
		return customErr.HandleDBError(err)
	}
	return nil
}
