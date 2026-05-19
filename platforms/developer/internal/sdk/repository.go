package sdk

import (
	"context"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, sdk *SDK) error
	GetByID(ctx context.Context, id string) (*SDK, error)
	ListByLanguage(ctx context.Context, language string) ([]*SDK, error)
	List(ctx context.Context) ([]*SDK, error)
	Update(ctx context.Context, sdk *SDK) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	sdks map[string]*SDK
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sdks: make(map[string]*SDK),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, sdk *SDK) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sdks[sdk.ID] = sdk
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*SDK, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sdks[id]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*SDK, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*SDK
	for _, s := range r.sdks {
		result = append(result, s)
	}
	return result, nil
}

func (r *InMemoryRepository) ListByLanguage(ctx context.Context, language string) ([]*SDK, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*SDK
	for _, s := range r.sdks {
		if s.Language == language {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, sdk *SDK) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sdks[sdk.ID] = sdk
	return nil
}
