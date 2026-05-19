package recommender

import (
	"context"
	"sync"
)

type Repository interface {
	GetRecommendations(ctx context.Context, ctxs RecommendationContext) ([]ProductRecommendation, error)
	StoreRecommendations(ctx context.Context, recs []ProductRecommendation) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data []ProductRecommendation
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make([]ProductRecommendation, 0),
	}
}

func (r *InMemoryRepository) GetRecommendations(ctx context.Context, recCtx RecommendationContext) ([]ProductRecommendation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ProductRecommendation, len(r.data))
	copy(result, r.data)
	return result, nil
}

func (r *InMemoryRepository) StoreRecommendations(ctx context.Context, recs []ProductRecommendation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = recs
	return nil
}
