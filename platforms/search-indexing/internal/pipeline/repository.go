package pipeline

import (
	"context"
	"sync"
)

type Repository interface {
	CreatePipeline(ctx context.Context, p *Pipeline) error
	GetPipeline(ctx context.Context, id string) (*Pipeline, error)
	ListPipelines(ctx context.Context) ([]*Pipeline, error)
	UpdatePipeline(ctx context.Context, p *Pipeline) error
	DeletePipeline(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	pipelines map[string]*Pipeline
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		pipelines: make(map[string]*Pipeline),
	}
}

func (r *InMemoryRepository) CreatePipeline(_ context.Context, p *Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.pipelines[p.ID]; ok {
		return ErrPipelineExists
	}
	r.pipelines[p.ID] = p
	return nil
}

func (r *InMemoryRepository) GetPipeline(_ context.Context, id string) (*Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pipelines[id]
	if !ok {
		return nil, ErrPipelineNotFound
	}
	return p, nil
}

func (r *InMemoryRepository) ListPipelines(_ context.Context) ([]*Pipeline, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Pipeline, 0, len(r.pipelines))
	for _, p := range r.pipelines {
		list = append(list, p)
	}
	return list, nil
}

func (r *InMemoryRepository) UpdatePipeline(_ context.Context, p *Pipeline) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.pipelines[p.ID]; !ok {
		return ErrPipelineNotFound
	}
	r.pipelines[p.ID] = p
	return nil
}

func (r *InMemoryRepository) DeletePipeline(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.pipelines[id]; !ok {
		return ErrPipelineNotFound
	}
	delete(r.pipelines, id)
	return nil
}
