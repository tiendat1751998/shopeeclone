package content

import (
	"context"
	"sync"
)

type Repository interface {
	GetProductFeatures(ctx context.Context, productID string) (*ProductFeatures, error)
	GetAllProductFeatures(ctx context.Context) ([]ProductFeatures, error)
	StoreProductFeatures(ctx context.Context, features *ProductFeatures) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	features map[string]*ProductFeatures
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		features: make(map[string]*ProductFeatures),
	}
}

func (r *InMemoryRepository) GetProductFeatures(ctx context.Context, productID string) (*ProductFeatures, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.features[productID]
	if !ok {
		return nil, nil
	}
	fCopy := *f
	return &fCopy, nil
}

func (r *InMemoryRepository) GetAllProductFeatures(ctx context.Context) ([]ProductFeatures, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]ProductFeatures, 0, len(r.features))
	for _, f := range r.features {
		result = append(result, *f)
	}
	return result, nil
}

func (r *InMemoryRepository) StoreProductFeatures(ctx context.Context, features *ProductFeatures) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	fCopy := *features
	r.features[features.ProductID] = &fCopy
	return nil
}
