package dispute

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, d *Dispute) error
	Update(ctx context.Context, d *Dispute) error
	GetByID(ctx context.Context, id string) (*Dispute, error)
	List(ctx context.Context, offset, limit int) ([]*Dispute, int64, error)
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Dispute
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{data: make(map[string]*Dispute)}
}

func (r *InMemoryRepository) Create(ctx context.Context, d *Dispute) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	d.OpenedAt = now
	r.data[d.ID] = d
	return nil
}

func (r *InMemoryRepository) Update(ctx context.Context, d *Dispute) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[d.ID]; !ok {
		return ErrDisputeNotFound
	}
	r.data[d.ID] = d
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Dispute, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.data[id]
	if !ok {
		return nil, ErrDisputeNotFound
	}
	return d, nil
}

func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]*Dispute, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := int64(len(r.data))
	items := make([]*Dispute, 0, limit)
	i := 0
	for _, d := range r.data {
		if i >= offset && len(items) < limit {
			items = append(items, d)
		}
		i++
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].OpenedAt > items[j].OpenedAt
	})
	return items, total, nil
}
