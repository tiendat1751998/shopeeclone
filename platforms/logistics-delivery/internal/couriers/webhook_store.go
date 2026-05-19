package couriers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisWebhookStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisWebhookStore(client *redis.Client, ttl time.Duration) *RedisWebhookStore {
	return &RedisWebhookStore{client: client, ttl: ttl}
}

func (s *RedisWebhookStore) IsDuplicate(ctx context.Context, eventID string) (bool, error) {
	exists, err := s.client.Exists(ctx, "webhook:"+eventID).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (s *RedisWebhookStore) MarkProcessed(ctx context.Context, eventID string) error {
	return s.client.Set(ctx, "webhook:"+eventID, "1", s.ttl).Err()
}
