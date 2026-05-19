package trending

import (
	"context"
	"sync"
	"time"
)

type Interaction struct {
	ProductID string    `json:"product_id"`
	Timestamp time.Time `json:"timestamp"`
	Weight    float64   `json:"weight"`
}

type Repository interface {
	RecordInteraction(ctx context.Context, productID string, weight float64) error
	GetWindowInteractions(ctx context.Context, since time.Time) ([]Interaction, error)
	GetAllProductIDs(ctx context.Context) ([]string, error)
}

type InMemoryRepository struct {
	mu           sync.RWMutex
	interactions []Interaction
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		interactions: make([]Interaction, 0),
	}
}

func (r *InMemoryRepository) RecordInteraction(ctx context.Context, productID string, weight float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.interactions = append(r.interactions, Interaction{
		ProductID: productID,
		Timestamp: time.Now(),
		Weight:    weight,
	})
	return nil
}

func (r *InMemoryRepository) GetWindowInteractions(ctx context.Context, since time.Time) ([]Interaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Interaction
	for _, inter := range r.interactions {
		if inter.Timestamp.After(since) {
			result = append(result, inter)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) GetAllProductIDs(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	seen := make(map[string]bool)
	var ids []string
	for _, inter := range r.interactions {
		if !seen[inter.ProductID] {
			seen[inter.ProductID] = true
			ids = append(ids, inter.ProductID)
		}
	}
	return ids, nil
}
