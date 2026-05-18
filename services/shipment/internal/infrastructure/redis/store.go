package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/services/shipment/internal/config"
)

type Store struct {
	client *redis.Client
	cfg    config.RedisConfig
}

func NewStore(client *redis.Client, cfg config.RedisConfig) *Store { return &Store{client: client, cfg: cfg} }
func (s *Store) IsAvailable() bool { return s.client != nil }

func (s *Store) CheckIdempotencyKey(ctx context.Context, key string) (string, error) {
	if s.client == nil { return "", nil }
	return s.client.Get(ctx, fmt.Sprintf("shipment:idempotency:%s", key)).Result()
}

func (s *Store) StoreIdempotencyKey(ctx context.Context, key, shipmentID string, ttl time.Duration) error {
	if s.client == nil { return nil }
	return s.client.Set(ctx, fmt.Sprintf("shipment:idempotency:%s", key), shipmentID, ttl).Err()
}

func (s *Store) CheckWebhookReplay(ctx context.Context, key string) (bool, error) {
	if s.client == nil { return false, nil }
	exists, err := s.client.Exists(ctx, fmt.Sprintf("shipment:webhook:%s", key)).Result()
	return exists > 0, err
}

func (s *Store) MarkWebhookProcessed(ctx context.Context, key string, ttl time.Duration) error {
	if s.client == nil { return nil }
	return s.client.Set(ctx, fmt.Sprintf("shipment:webhook:%s", key), "1", ttl).Err()
}

func (s *Store) AcquireShipmentLock(ctx context.Context, shipmentID string, ttl time.Duration) (bool, error) {
	if s.client == nil { return true, nil }
	return s.client.SetNX(ctx, fmt.Sprintf("lock:shipment:%s", shipmentID), "1", ttl).Result()
}

func (s *Store) ReleaseShipmentLock(ctx context.Context, shipmentID string) error {
	if s.client == nil { return nil }
	return s.client.Del(ctx, fmt.Sprintf("lock:shipment:%s", shipmentID)).Err()
}
