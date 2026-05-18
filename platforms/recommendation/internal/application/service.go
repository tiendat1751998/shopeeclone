package application
import ("context"; "encoding/json"; "fmt"; "time"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/domain"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/infrastructure/redis"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/metrics"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.opentelemetry.io/otel/attribute"; "go.uber.org/zap")

type RecommendationService struct { recRepo domain.RecommendationRepository; eventRepo domain.UserEventRepository; features domain.FeatureStore; redis *redis.Store; publisher EventPublisher }
type EventPublisher interface { Publish(ctx context.Context, eventType string, payload interface{}) error }

func NewRecommendationService(rr domain.RecommendationRepository, er domain.UserEventRepository, fs domain.FeatureStore, rs *redis.Store, pub EventPublisher) *RecommendationService {
	return &RecommendationService{recRepo: rr, eventRepo: er, features: fs, redis: rs, publisher: pub}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, req domain.RecommendationRequest) (*domain.RecommendationResponse, error) {
	ctx, span := otel.Tracer("shopee-recommendation").Start(ctx, "rec.get"); defer span.End()
	span.SetAttributes(attribute.String("user_id", req.UserID), attribute.String("type", req.Context))

	if data, err := s.redis.GetCachedRecommendations(ctx, req.UserID, req.Context); err == nil && len(data) > 0 {
		metrics.CacheHitsTotal.WithLabelValues("recommendation", "redis").Inc()
		var recs []domain.Recommendation; json.Unmarshal(data, &recs)
		return &domain.RecommendationResponse{Recommendations: recs, TookMs: 2}, nil
	}
	metrics.CacheMissesTotal.WithLabelValues("recommendation", "redis").Inc()

	limit := req.Limit; if limit == 0 { limit = 20 }
	recs, err := s.recRepo.GetRecommendations(ctx, req.UserID, req.Context, limit)
	if err != nil { metrics.RecErrors.Inc(); return nil, domain.ErrRecFailed }

	data, _ := json.Marshal(recs); s.redis.SetCachedRecommendations(ctx, req.UserID, req.Context, data, 15*time.Minute)
	metrics.RecRequestsTotal.Inc()
	return &domain.RecommendationResponse{Recommendations: recs, TookMs: 10}, nil
}

func (s *RecommendationService) TrackEvent(ctx context.Context, event *domain.UserEvent) error {
	if err := s.eventRepo.TrackEvent(ctx, event); err != nil { return err }
	if event.EventType == "product_view" {
		s.redis.IncrementProductView(ctx, event.ProductID)
	}
	metrics.EventsTracked.Inc()
	return nil
}

func (s *RecommendationService) GetTrending(ctx context.Context, limit int) ([]domain.Recommendation, error) {
	productIDs, err := s.redis.GetTrendingProducts(ctx, limit)
	if err != nil { return nil, err }
	var recs []domain.Recommendation
	for _, pid := range productIDs {
		recs = append(recs, domain.Recommendation{ProductID: pid, Type: domain.RecTypeTrending, Score: 1.0})
	}
	return recs, nil
}

func (s *RecommendationService) GetSimilarProducts(ctx context.Context, productID string, limit int) ([]domain.Recommendation, error) {
	return s.recRepo.GetSimilarProducts(ctx, productID, limit)
}
