package redis
import ("context"; "fmt"; "time"; "github.com/redis/go-redis/v9"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/config")
type Store struct { rdb *redis.Client; cfg config.RedisConfig }
func NewStore(rdb *redis.Client, cfg config.RedisConfig) *Store { return &Store{rdb: rdb, cfg: cfg} }
func (s *Store) Ping(ctx context.Context) error { return s.rdb.Ping(ctx).Err() }
func (s *Store) Close() error { return s.rdb.Close() }
func (s *Store) GetCachedRecommendations(ctx context.Context, userID, recType string) ([]byte, error) {
	return s.rdb.Get(ctx, fmt.Sprintf("rec:%s:%s", userID, recType)).Bytes()
}
func (s *Store) SetCachedRecommendations(ctx context.Context, userID, recType string, data []byte, ttl time.Duration) error {
	return s.rdb.Set(ctx, fmt.Sprintf("rec:%s:%s", userID, recType), data, ttl).Err()
}
func (s *Store) IncrementProductView(ctx context.Context, productID string) error {
	return s.rdb.ZIncrBy(ctx, "trending:products", 1, productID).Err()
}
func (s *Store) GetTrendingProducts(ctx context.Context, limit int) ([]string, error) {
	return s.rdb.ZRevRange(ctx, "trending:products", 0, int64(limit-1)).Result()
}
func (s *Store) SetUserFeatures(ctx context.Context, userID string, data []byte, ttl time.Duration) error {
	return s.rdb.Set(ctx, fmt.Sprintf("features:user:%s", userID), data, ttl).Err()
}
func (s *Store) GetUserFeatures(ctx context.Context, userID string) ([]byte, error) {
	return s.rdb.Get(ctx, fmt.Sprintf("features:user:%s", userID)).Bytes()
}
