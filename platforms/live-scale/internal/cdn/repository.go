package cdn

import (
	"context"
	"sync"
)

type Repository interface {
	SaveEndpoint(ctx context.Context, ep *CDNEndpoint) error
	GetEndpoint(ctx context.Context, id string) (*CDNEndpoint, error)
	ListEndpoints(ctx context.Context) ([]*CDNEndpoint, error)
	CreatePurgeRequest(ctx context.Context, req *CDNPurgeRequest) error
	ListPurgeRequests(ctx context.Context) ([]*CDNPurgeRequest, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	endpoints map[string]*CDNEndpoint
	purges    map[string]*CDNPurgeRequest
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		endpoints: make(map[string]*CDNEndpoint),
		purges:    make(map[string]*CDNPurgeRequest),
	}
}

func (r *InMemoryRepository) SaveEndpoint(_ context.Context, ep *CDNEndpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.endpoints[ep.ID] = ep
	return nil
}

func (r *InMemoryRepository) GetEndpoint(_ context.Context, id string) (*CDNEndpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ep, ok := r.endpoints[id]
	if !ok {
		return nil, ErrEndpointNotFound
	}
	return ep, nil
}

func (r *InMemoryRepository) ListEndpoints(_ context.Context) ([]*CDNEndpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*CDNEndpoint, 0, len(r.endpoints))
	for _, ep := range r.endpoints {
		result = append(result, ep)
	}
	return result, nil
}

func (r *InMemoryRepository) CreatePurgeRequest(_ context.Context, req *CDNPurgeRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.purges[req.ID] = req
	return nil
}

func (r *InMemoryRepository) ListPurgeRequests(_ context.Context) ([]*CDNPurgeRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*CDNPurgeRequest, 0, len(r.purges))
	for _, p := range r.purges {
		result = append(result, p)
	}
	return result, nil
}
