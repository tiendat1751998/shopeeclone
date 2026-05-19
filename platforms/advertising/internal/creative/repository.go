package creative

import (
	"context"
	"sync"
)

type Repository interface {
	Create(ctx context.Context, c *Creative) (*Creative, error)
	Get(ctx context.Context, id string) (*Creative, error)
	GetByCampaign(ctx context.Context, campaignID string) ([]*Creative, error)
	List(ctx context.Context, status CreativeStatus) ([]*Creative, error)
	Update(ctx context.Context, c *Creative) (*Creative, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	creatives map[string]*Creative
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		creatives: make(map[string]*Creative),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, c *Creative) (*Creative, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creatives[c.ID] = c
	return c, nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*Creative, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.creatives[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryRepository) GetByCampaign(ctx context.Context, campaignID string) ([]*Creative, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Creative
	for _, c := range r.creatives {
		if c.CampaignID == campaignID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) List(ctx context.Context, status CreativeStatus) ([]*Creative, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Creative
	for _, c := range r.creatives {
		if status == "" || c.Status == status {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, c *Creative) (*Creative, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creatives[c.ID] = c
	return c, nil
}
