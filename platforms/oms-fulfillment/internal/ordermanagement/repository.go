package ordermanagement

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Repository interface {
	Create(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	Update(ctx context.Context, o *Order) error
	List(ctx context.Context, filter OrderFilter) ([]*Order, int64, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{orders: make(map[string]*Order)}
}

func (r *InMemoryRepository) Create(_ context.Context, o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID] = o
	return nil
}

func (r *InMemoryRepository) GetByID(_ context.Context, id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return o, nil
}

func (r *InMemoryRepository) Update(_ context.Context, o *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[o.ID]; !ok {
		return ErrOrderNotFound
	}
	r.orders[o.ID] = o
	return nil
}

func (r *InMemoryRepository) List(_ context.Context, filter OrderFilter) ([]*Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Order
	for _, o := range r.orders {
		if filter.Status != "" && o.Status != filter.Status {
			continue
		}
		if filter.UserID != "" && o.UserID != filter.UserID {
			continue
		}
		if filter.From != nil && o.CreatedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && o.CreatedAt.After(*filter.To) {
			continue
		}
		result = append(result, o)
	}
	start := filter.Offset
	if start > len(result) {
		start = len(result)
	}
	end := start + filter.Limit
	if filter.Limit <= 0 {
		end = len(result)
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], int64(len(result)), nil
}

// Ensure interface compliance
var _ Repository = (*InMemoryRepository)(nil)

// Helper for generating IDs
func GenerateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
