package redis
import ("context"; "fmt"; "time"; "github.com/redis/go-redis/v9"; "github.com/shopee-clone/shopee/platforms/search/internal/config")
type Store struct { rdb *redis.Client; cfg config.RedisConfig }
func NewStore(rdb *redis.Client, cfg config.RedisConfig) *Store { return &Store{rdb: rdb, cfg: cfg} }
func (s *Store) Ping(ctx context.Context) error { return s.rdb.Ping(ctx).Err() }
func (s *Store) Close() error { return s.rdb.Close() }
func (s *Store) GetCachedSearch(ctx context.Context, key string) ([]byte, error) { return s.rdb.Get(ctx, "search:"+key).Bytes() }
func (s *Store) SetCachedSearch(ctx context.Context, key string, data []byte, ttl time.Duration) error { return s.rdb.Set(ctx, "search:"+key, data, ttl).Err() }
func (s *Store) GetCachedAutocomplete(ctx context.Context, prefix string) ([]byte, error) { return s.rdb.Get(ctx, "autocomplete:"+prefix).Bytes() }
func (s *Store) SetCachedAutocomplete(ctx context.Context, prefix string, data []byte, ttl time.Duration) error { return s.rdb.Set(ctx, "autocomplete:"+prefix, data, ttl).Err() }
func (s *Store) IncrementQueryCounter(ctx context.Context, query string) error { return s.rdb.ZIncrBy(ctx, "trending_queries", 1, query).Err() }
func (s *Store) GetTrendingQueries(ctx context.Context, limit int) ([]string, error) { return s.rdb.ZRevRange(ctx, "trending_queries", 0, int64(limit-1)).Result() }
