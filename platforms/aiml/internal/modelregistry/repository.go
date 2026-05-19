package modelregistry

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Store(ctx context.Context, model *Model) error
	Get(ctx context.Context, id string) (*Model, error)
	List(ctx context.Context) ([]*Model, error)
	ListByStage(ctx context.Context, stage Stage) ([]*Model, error)
	Update(ctx context.Context, model *Model) error
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	models map[string]*Model
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		models: make(map[string]*Model),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, model *Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.models[model.ID]; ok {
		return ErrModelExists
	}
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	r.models[model.ID] = model
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.models[id]
	if !ok {
		return nil, ErrModelNotFound
	}
	return m, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Model, 0, len(r.models))
	for _, m := range r.models {
		result = append(result, m)
	}
	return result, nil
}

func (r *InMemoryRepository) ListByStage(ctx context.Context, stage Stage) ([]*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Model
	for _, m := range r.models {
		if m.Status == stage {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, model *Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.models[model.ID]; !ok {
		return ErrModelNotFound
	}
	r.models[model.ID] = model
	return nil
}
