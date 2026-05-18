package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/services/inventory/internal/config"
)

type Store struct {
	client *redis.Client
	cfg    config.RedisConfig
}

func NewStore(client *redis.Client, cfg config.RedisConfig) *Store { return &Store{client: client, cfg: cfg} }
func (s *Store) IsAvailable() bool { return s.client != nil }

func (s *Store) CheckIdempotencyKey(ctx context.Context, key string) (string, error) {
	if s.client == nil { return "", nil }
	return s.client.Get(ctx, fmt.Sprintf("inventory:idempotency:%s", key)).Result()
}

func (s *Store) StoreIdempotencyKey(ctx context.Context, key, reservationID string, ttl time.Duration) error {
	if s.client == nil { return nil }
	return s.client.Set(ctx, fmt.Sprintf("inventory:idempotency:%s", key), reservationID, ttl).Err()
}

func (s *Store) AcquireStockLock(ctx context.Context, skuID string, ttl time.Duration) (bool, error) {
	if s.client == nil { return true, nil }
	return s.client.SetNX(ctx, fmt.Sprintf("lock:inventory:%s", skuID), "1", ttl).Result()
}

func (s *Store) ReleaseStockLock(ctx context.Context, skuID string) error {
	if s.client == nil { return nil }
	return s.client.Del(ctx, fmt.Sprintf("lock:inventory:%s", skuID)).Err()
}

func (s *Store) IncrementCounter(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	if s.client == nil { return 0, nil }
	pipe := s.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil { return 0, err }
	return incr.Val(), nil
}

func (s *Store) CacheStock(ctx context.Context, skuID string, quantity, reserved int, ttl time.Duration) error {
	if s.client == nil { return nil }
	return s.client.HMSet(ctx, fmt.Sprintf("stock:%s", skuID), "quantity", quantity, "reserved", reserved, "available", quantity-reserved).Err()
}

func (s *Store) GetCachedStock(ctx context.Context, skuID string) (int, int, error) {
	if s.client == nil { return 0, 0, fmt.Errorf("redis not available") }
	vals, err := s.client.HMGet(ctx, fmt.Sprintf("stock:%s", skuID), "quantity", "reserved").Result()
	if err != nil { return 0, 0, err }
	qty, _ := vals[0].(int64)
	res, _ := vals[1].(int64)
	return int(qty), int(res), nil
}
