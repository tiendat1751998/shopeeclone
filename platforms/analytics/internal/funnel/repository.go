package funnel

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreDefinition(ctx context.Context, def *FunnelDefinition) error
	GetDefinition(ctx context.Context, id string) (*FunnelDefinition, error)
	StoreResult(ctx context.Context, result *FunnelResult) error
	GetResult(ctx context.Context, id string) (*FunnelResult, error)
	ListResults(ctx context.Context, offset, limit int) ([]*FunnelResult, int, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	definitions map[string]*FunnelDefinition
	results     map[string]*FunnelResult
	resultsList []*FunnelResult
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		definitions: make(map[string]*FunnelDefinition),
		results:     make(map[string]*FunnelResult),
		resultsList: make([]*FunnelResult, 0),
	}
}

func (r *InMemoryRepository) StoreDefinition(ctx context.Context, def *FunnelDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if def.CreatedAt.IsZero() {
		def.CreatedAt = time.Now()
	}
	r.definitions[def.ID] = def
	return nil
}

func (r *InMemoryRepository) GetDefinition(ctx context.Context, id string) (*FunnelDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	def, ok := r.definitions[id]
	if !ok {
		return nil, nil
	}
	return def, nil
}

func (r *InMemoryRepository) StoreResult(ctx context.Context, result *FunnelResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if result.AnalyzedAt.IsZero() {
		result.AnalyzedAt = time.Now()
	}
	r.results[result.ID] = result
	r.resultsList = append(r.resultsList, result)
	return nil
}

func (r *InMemoryRepository) GetResult(ctx context.Context, id string) (*FunnelResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result, ok := r.results[id]
	if !ok {
		return nil, nil
	}
	return result, nil
}

func (r *InMemoryRepository) ListResults(ctx context.Context, offset, limit int) ([]*FunnelResult, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.resultsList)
	if offset >= total {
		return []*FunnelResult{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return r.resultsList[offset:end], total, nil
}
