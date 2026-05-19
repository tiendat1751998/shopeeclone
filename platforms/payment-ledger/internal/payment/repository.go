package payment

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, p *Payment) error
	Update(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	List(ctx context.Context, offset, limit int) ([]*Payment, int64, error)
	GetByOrder(ctx context.Context, orderID string) ([]*Payment, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	data  map[string]*Payment
	order map[string][]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data:  make(map[string]*Payment),
		order: make(map[string][]string),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, p *Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	p.CreatedAt = now
	p.UpdatedAt = now
	r.data[p.ID] = p
	r.order[p.OrderID] = append(r.order[p.OrderID], p.ID)
	return nil
}

func (r *InMemoryRepository) Update(ctx context.Context, p *Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.data[p.ID]
	if !ok {
		return ErrPaymentNotFound
	}
	p.CreatedAt = existing.CreatedAt
	p.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.data[p.ID] = p
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, ErrPaymentNotFound
	}
	return p, nil
}

func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]*Payment, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := int64(len(r.data))
	items := make([]*Payment, 0, limit)
	i := 0
	for _, p := range r.data {
		if i >= offset && len(items) < limit {
			items = append(items, p)
		}
		i++
	}
	return items, total, nil
}

func (r *InMemoryRepository) GetByOrder(ctx context.Context, orderID string) ([]*Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.order[orderID]
	if !ok {
		return nil, nil
	}
	result := make([]*Payment, 0, len(ids))
	for _, id := range ids {
		if p, ok := r.data[id]; ok {
			result = append(result, p)
		}
	}
	return result, nil
}
