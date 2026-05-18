package application
import ("context"; "encoding/json"; "fmt"; "time"; "github.com/shopee-clone/shopee/platforms/search/internal/domain"; "github.com/shopee-clone/shopee/platforms/search/internal/infrastructure/elasticsearch"; "github.com/shopee-clone/shopee/platforms/search/internal/infrastructure/redis"; "github.com/shopee-clone/shopee/platforms/search/internal/metrics"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.opentelemetry.io/otel/attribute"; "go.uber.org/zap")

type SearchService struct { esClient *elasticsearch.Client; redis *redis.Store; publisher EventPublisher }
type EventPublisher interface { Publish(ctx context.Context, eventType string, payload interface{}) error }

func NewSearchService(es *elasticsearch.Client, rs *redis.Store, pub EventPublisher) *SearchService {
	return &SearchService{esClient: es, redis: rs, publisher: pub}
}

func (s *SearchService) Search(ctx context.Context, query domain.SearchQuery) (*domain.SearchResult, error) {
	ctx, span := otel.Tracer("shopee-search").Start(ctx, "search.query"); defer span.End()
	span.SetAttributes(attribute.String("query", query.Query), attribute.Int("page", query.Page))

	// Check cache
	cacheKey := fmt.Sprintf("%s:%s:%d:%d", query.Query, query.CategoryID, query.Page, query.Limit)
	if data, err := s.redis.GetCachedSearch(ctx, cacheKey); err == nil && len(data) > 0 {
		metrics.CacheHitsTotal.WithLabelValues("search", "redis").Inc()
		var result domain.SearchResult; json.Unmarshal(data, &result); return &result, nil
	}
	metrics.CacheMissesTotal.WithLabelValues("search", "redis").Inc()

	// Search ES
	result, err := s.esClient.Search(ctx, query)
	if err != nil { metrics.SearchErrors.Inc(); return nil, domain.ErrSearchFailed }

	// Cache result
	data, _ := json.Marshal(result); s.redis.SetCachedSearch(ctx, cacheKey, data, 5*time.Minute)

	// Record analytics
	s.redis.IncrementQueryCounter(ctx, query.Query)
	metrics.SearchQueriesTotal.Inc()

	return result, nil
}

func (s *SearchService) Autocomplete(ctx context.Context, prefix string, limit int) (*domain.AutocompleteResult, error) {
	if data, err := s.redis.GetCachedAutocomplete(ctx, prefix); err == nil && len(data) > 0 {
		var result domain.AutocompleteResult; json.Unmarshal(data, &result); return &result, nil
	}
	result, err := s.esClient.Autocomplete(ctx, prefix, limit)
	if err != nil { return nil, err }
	data, _ := json.Marshal(result); s.redis.SetCachedAutocomplete(ctx, prefix, data, 10*time.Minute)
	metrics.AutocompleteRequestsTotal.Inc()
	return result, nil
}

func (s *SearchService) IndexProduct(ctx context.Context, doc *domain.IndexDocument) error {
	if err := s.esClient.Index(ctx, doc); err != nil { return err }
	s.redis.SetCachedSearch(ctx, doc.ID, nil, 0) // invalidate
	metrics.DocumentsIndexed.Inc()
	if s.publisher != nil { s.publisher.Publish(ctx, "indexing.triggered", doc) }
	return nil
}

func (s *SearchService) BulkIndex(ctx context.Context, docs []*domain.IndexDocument) error {
	return s.esClient.BulkIndex(ctx, docs)
}

func (s *SearchService) GetTrendingQueries(ctx context.Context, limit int) ([]string, error) {
	return s.redis.GetTrendingQueries(ctx, limit)
}
