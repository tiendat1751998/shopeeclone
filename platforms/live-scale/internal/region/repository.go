package region

import (
	"context"
	"sync"
)

type Repository interface {
	SaveRegion(ctx context.Context, region *Region) error
	GetRegion(ctx context.Context, code string) (*Region, error)
	ListRegions(ctx context.Context) ([]*Region, error)
	SaveLatency(ctx context.Context, latency *LatencyMap) error
	ListLatencies(ctx context.Context) ([]*LatencyMap, error)
	SaveRoutingRule(ctx context.Context, rule *GeoRoutingRule) error
	ListRoutingRules(ctx context.Context) ([]*GeoRoutingRule, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	regions  map[string]*Region
	latencies []*LatencyMap
	rules    map[string]*GeoRoutingRule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		regions:   make(map[string]*Region),
		latencies: make([]*LatencyMap, 0),
		rules:     make(map[string]*GeoRoutingRule),
	}
}

func (r *InMemoryRepository) SaveRegion(_ context.Context, region *Region) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.regions[region.Code] = region
	return nil
}

func (r *InMemoryRepository) GetRegion(_ context.Context, code string) (*Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reg, ok := r.regions[code]
	if !ok {
		return nil, ErrRegionNotFound
	}
	return reg, nil
}

func (r *InMemoryRepository) ListRegions(_ context.Context) ([]*Region, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Region, 0, len(r.regions))
	for _, reg := range r.regions {
		result = append(result, reg)
	}
	return result, nil
}

func (r *InMemoryRepository) SaveLatency(_ context.Context, latency *LatencyMap) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, lm := range r.latencies {
		if lm.FromRegion == latency.FromRegion && lm.ToRegion == latency.ToRegion {
			r.latencies[i] = latency
			return nil
		}
	}
	r.latencies = append(r.latencies, latency)
	return nil
}

func (r *InMemoryRepository) ListLatencies(_ context.Context) ([]*LatencyMap, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*LatencyMap, len(r.latencies))
	copy(result, r.latencies)
	return result, nil
}

func (r *InMemoryRepository) SaveRoutingRule(_ context.Context, rule *GeoRoutingRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *InMemoryRepository) ListRoutingRules(_ context.Context) ([]*GeoRoutingRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*GeoRoutingRule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}
