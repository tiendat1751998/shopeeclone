package redis

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

type Store struct{ client *redis.Client }

func NewStore(client *redis.Client) *Store { return &Store{client: client} }
func (s *Store) Client() *redis.Client { return s.client }

func (s *Store) AcquireIdempotencyLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	key = fmt.Sprintf("idempotency:%s", key)
	ok, err := s.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *Store) IsProcessed(ctx context.Context, key string) (bool, error) {
	key = fmt.Sprintf("idempotency:%s", key)
	_, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) MarkProcessed(ctx context.Context, key string, ttl time.Duration) error {
	key = fmt.Sprintf("idempotency:%s", key)
	return s.client.Set(ctx, key, "1", ttl).Err()
}

func (s *Store) CacheBalance(ctx context.Context, accountID string, balance int64) error {
	return s.client.Set(ctx, fmt.Sprintf("balance:%s", accountID), balance, 30*time.Second).Err()
}

func (s *Store) GetCachedBalance(ctx context.Context, accountID string) (int64, error) {
	return s.client.Get(ctx, fmt.Sprintf("balance:%s", accountID)).Int64()
}

func (s *Store) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	key = fmt.Sprintf("ratelimit:%s", key)
	current, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if current == 1 {
		s.client.Expire(ctx, key, window)
	}
	return current <= int64(limit), nil
}
