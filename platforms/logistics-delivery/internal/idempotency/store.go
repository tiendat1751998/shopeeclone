package idempotency

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ttl    time.Duration
}

type KeyInfo struct {
	Key       string    `json:"key"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}

func NewStore(client *redis.Client, ttl time.Duration) *Store {
	return &Store{client: client, ttl: ttl}
}

func (s *Store) IsDuplicate(ctx context.Context, key string) (bool, error) {
	exists, err := s.client.Exists(ctx, "idempotency:"+key).Result()
	if err != nil {
		return false, fmt.Errorf("check idempotency: %w", err)
	}
	return exists > 0, nil
}

func (s *Store) MarkProcessed(ctx context.Context, key string) error {
	info := KeyInfo{
		Key:       key,
		Processed: true,
		CreatedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal idempotency: %w", err)
	}
	return s.client.Set(ctx, "idempotency:"+key, data, s.ttl).Err()
}

func (s *Store) TryLock(ctx context.Context, key string) (bool, error) {
	ok, err := s.client.SetNX(ctx, "lock:"+key, "1", 30*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("lock: %w", err)
	}
	return ok, nil
}

func (s *Store) ReleaseLock(ctx context.Context, key string) error {
	return s.client.Del(ctx, "lock:"+key).Err()
}
