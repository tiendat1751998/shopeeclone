package domain
import "context"

type RecommendationRepository interface {
	GetRecommendations(ctx context.Context, userID, recType string, limit int) ([]Recommendation, error)
	SaveRecommendations(ctx context.Context, recs []Recommendation) error
	GetSimilarProducts(ctx context.Context, productID string, limit int) ([]Recommendation, error)
	GetTrending(ctx context.Context, limit int) ([]Recommendation, error)
}

type UserEventRepository interface {
	TrackEvent(ctx context.Context, event *UserEvent) error
	GetUserEvents(ctx context.Context, userID string, limit int) ([]UserEvent, error)
}

type FeatureStore interface {
	GetUserFeatures(ctx context.Context, userID string) (*FeatureVector, error)
	GetProductFeatures(ctx context.Context, productID string) (*FeatureVector, error)
	UpdateFeatures(ctx context.Context, vector *FeatureVector) error
}
