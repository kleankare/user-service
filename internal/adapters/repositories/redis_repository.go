package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kleankare/user-service/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) ports.CacheRepository {
	return &redisRepository{
		client: client,
	}
}

func (repo *redisRepository) Get(key string, value interface{}) error {
	data, err := repo.client.Get(context.Background(), key).Result()

	if err == redis.Nil {
		return fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return fmt.Errorf("failed to get value for key %q: %v", key, err)
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %q: %v", key, err)
	}

	return nil
}

func (repo *redisRepository) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %v", key, err)
	}
	if err := repo.client.Set(context.Background(), key, data, duration).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %q: %v", key, err)
	}

	return nil
}

func (repo *redisRepository) Delete(key string) error {
	if err := repo.client.Del(context.Background(), key).Err(); err != nil {
		return fmt.Errorf("failed to delete value for key %q: %v", key, err)
	}
	return nil
}
