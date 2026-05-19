package campaign

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *Campaign) error
	GetByID(ctx context.Context, id string) (*Campaign, error)
	List(ctx context.Context) ([]*Campaign, error)
	Update(ctx context.Context, c *Campaign) error
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	campaigns map[string]*Campaign
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		campaigns: make(map[string]*Campaign),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, c *Campaign) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.Status == "" {
		c.Status = StatusDraft
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[c.ID] = c
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.campaigns[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Campaign
	for _, c := range r.campaigns {
		result = append(result, c)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, c *Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.campaigns[c.ID]
	if !ok {
		return nil
	}
	c.CreatedAt = existing.CreatedAt
	c.UpdatedAt = time.Now()
	r.campaigns[c.ID] = c
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.campaigns, id)
	return nil
}
