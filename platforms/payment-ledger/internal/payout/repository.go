package payout

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PayoutRepository interface {
	Create(ctx context.Context, p *Payout) error
	Update(ctx context.Context, p *Payout) error
	GetByID(ctx context.Context, id string) (*Payout, error)
	List(ctx context.Context, offset, limit int) ([]*Payout, int64, error)
	GetByBatch(ctx context.Context, batchID string) ([]*Payout, error)
}

type BatchRepository interface {
	Create(ctx context.Context, b *PayoutBatch) error
	GetByID(ctx context.Context, id string) (*PayoutBatch, error)
	Update(ctx context.Context, b *PayoutBatch) error
}

type InMemoryPayoutRepo struct {
	mu    sync.RWMutex
	data  map[string]*Payout
	batch map[string][]string
}

func NewInMemoryPayoutRepo() *InMemoryPayoutRepo {
	return &InMemoryPayoutRepo{
		data:  make(map[string]*Payout),
		batch: make(map[string][]string),
	}
}

func (r *InMemoryPayoutRepo) Create(ctx context.Context, p *Payout) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	p.CreatedAt = now
	r.data[p.ID] = p
	if p.BatchID != "" {
		r.batch[p.BatchID] = append(r.batch[p.BatchID], p.ID)
	}
	return nil
}

func (r *InMemoryPayoutRepo) Update(ctx context.Context, p *Payout) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[p.ID]; !ok {
		return ErrPayoutNotFound
	}
	r.data[p.ID] = p
	return nil
}

func (r *InMemoryPayoutRepo) GetByID(ctx context.Context, id string) (*Payout, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, ErrPayoutNotFound
	}
	return p, nil
}

func (r *InMemoryPayoutRepo) List(ctx context.Context, offset, limit int) ([]*Payout, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := int64(len(r.data))
	items := make([]*Payout, 0, limit)
	i := 0
	for _, p := range r.data {
		if i >= offset && len(items) < limit {
			items = append(items, p)
		}
		i++
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, total, nil
}

func (r *InMemoryPayoutRepo) GetByBatch(ctx context.Context, batchID string) ([]*Payout, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids, ok := r.batch[batchID]
	if !ok {
		return nil, nil
	}
	result := make([]*Payout, 0, len(ids))
	for _, id := range ids {
		if p, ok := r.data[id]; ok {
			result = append(result, p)
		}
	}
	return result, nil
}

type InMemoryBatchRepo struct {
	mu   sync.RWMutex
	data map[string]*PayoutBatch
}

func NewInMemoryBatchRepo() *InMemoryBatchRepo {
	return &InMemoryBatchRepo{data: make(map[string]*PayoutBatch)}
}

func (r *InMemoryBatchRepo) Create(ctx context.Context, b *PayoutBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	b.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	r.data[b.ID] = b
	return nil
}

func (r *InMemoryBatchRepo) GetByID(ctx context.Context, id string) (*PayoutBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.data[id]
	if !ok {
		return nil, ErrPayoutNotFound
	}
	return b, nil
}

func (r *InMemoryBatchRepo) Update(ctx context.Context, b *PayoutBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[b.ID]; !ok {
		return ErrPayoutNotFound
	}
	r.data[b.ID] = b
	return nil
}
