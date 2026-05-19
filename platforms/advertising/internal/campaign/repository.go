package campaign

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Create(ctx context.Context, c *Campaign) error
	Get(ctx context.Context, id string) (*Campaign, error)
	List(ctx context.Context, status CampaignStatus) ([]*Campaign, error)
	Update(ctx context.Context, c *Campaign) (*Campaign, error)
	GetPerformance(ctx context.Context, id string) (*Performance, error)
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	campaigns  map[string]*Campaign
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		campaigns: make(map[string]*Campaign),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, c *Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[c.ID] = c
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.campaigns[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (r *InMemoryRepository) List(ctx context.Context, status CampaignStatus) ([]*Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Campaign
	for _, c := range r.campaigns {
		if status == "" || c.Status == status {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, c *Campaign) (*Campaign, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.campaigns[c.ID]; !ok {
		return nil, ErrNotFound
	}
	c.UpdatedAt = time.Now()
	r.campaigns[c.ID] = c
	return c, nil
}

func (r *InMemoryRepository) GetPerformance(ctx context.Context, id string) (*Performance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.campaigns[id]; !ok {
		return nil, ErrNotFound
	}
	return &Performance{
		CampaignID: id,
	}, nil
}
