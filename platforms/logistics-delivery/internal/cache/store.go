package cache

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

func NewStore(client *redis.Client, ttl time.Duration) *Store {
	return &Store{client: client, ttl: ttl}
}

func (s *Store) Get(ctx context.Context, key string, dest any) error {
	data, err := s.client.Get(ctx, "cache:"+key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (s *Store) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}
	return s.client.Set(ctx, "cache:"+key, data, s.ttl).Err()
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, "cache:"+key).Err()
}

func (s *Store) InvalidatePattern(ctx context.Context, pattern string) error {
	keys, err := s.client.Keys(ctx, "cache:"+pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}
	return nil
}
