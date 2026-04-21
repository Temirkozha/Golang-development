package repository

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	SetProcessing(ctx context.Context, key string, ttl time.Duration) (bool, error)
	SetCompleted(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (s *RedisStore) SetProcessing(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return s.client.SetNX(ctx, key, "processing", ttl).Result()
}

func (s *RedisStore) SetCompleted(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}