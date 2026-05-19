package cicd

import (
	"context"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, p *Pipeline) error
	GetByID(ctx context.Context, id string) (*Pipeline, error)
	List(ctx context.Context) ([]*Pipeline, error)
	Update(ctx context.Context, p *Pipeline) error
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	pipelines map[string]*Pipeline
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		pipelines: make(map[string]*Pipeline),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, p *Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pipelines[p.ID] = p
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pipelines[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Pipeline
	for _, p := range r.pipelines {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, p *Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pipelines[p.ID] = p
	return nil
}
