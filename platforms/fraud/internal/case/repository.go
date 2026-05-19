package fraudcase

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, c *FraudCase) error
	Get(ctx context.Context, id string) (*FraudCase, error)
	Update(ctx context.Context, c *FraudCase) error
	List(ctx context.Context, status CaseStatus, priority CasePriority, offset, limit int) ([]*FraudCase, int, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	cases  map[string]*FraudCase
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		cases: make(map[string]*FraudCase),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, c *FraudCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
	if c.Status == "" {
		c.Status = StatusOpen
	}
	r.cases[c.ID] = c
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (*FraudCase, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.cases[id]
	if !ok {
		return nil, ErrCaseNotFound
	}
	return c, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, c *FraudCase) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.UpdatedAt = time.Now()
	r.cases[c.ID] = c
	return nil
}

func (r *InMemoryRepository) List(ctx context.Context, status CaseStatus, priority CasePriority, offset, limit int) ([]*FraudCase, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*FraudCase
	for _, c := range r.cases {
		if status != "" && c.Status != status {
			continue
		}
		if priority != "" && c.Priority != priority {
			continue
		}
		filtered = append(filtered, c)
	}

	total := len(filtered)
	if offset >= total {
		return []*FraudCase{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}
