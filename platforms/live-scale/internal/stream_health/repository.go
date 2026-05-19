package stream_health

import (
	"context"
	"sync"
)

type Repository interface {
	SaveStreamHealth(ctx context.Context, health *StreamHealth) error
	GetStreamHealth(ctx context.Context, streamID string) (*StreamHealth, error)
	ListStreamHealth(ctx context.Context) ([]*StreamHealth, error)
	DeleteStreamHealth(ctx context.Context, streamID string) error
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	streams map[string]*StreamHealth
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		streams: make(map[string]*StreamHealth),
	}
}

func (r *InMemoryRepository) SaveStreamHealth(_ context.Context, health *StreamHealth) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streams[health.StreamID] = health
	return nil
}

func (r *InMemoryRepository) GetStreamHealth(_ context.Context, streamID string) (*StreamHealth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.streams[streamID]
	if !ok {
		return nil, ErrStreamNotFound
	}
	return h, nil
}

func (r *InMemoryRepository) ListStreamHealth(_ context.Context) ([]*StreamHealth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*StreamHealth, 0, len(r.streams))
	for _, h := range r.streams {
		result = append(result, h)
	}
	return result, nil
}

func (r *InMemoryRepository) DeleteStreamHealth(_ context.Context, streamID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.streams, streamID)
	return nil
}
