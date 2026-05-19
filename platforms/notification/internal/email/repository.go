package email

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, m *EmailMessage) error
	GetByID(ctx context.Context, id string) (*EmailMessage, error)
	UpdateStatus(ctx context.Context, id string, status EmailStatus) error
	ListByRecipient(ctx context.Context, email string, limit, offset int) ([]*EmailMessage, error)
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*EmailMessage
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*EmailMessage),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, m *EmailMessage) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.Status == "" {
		m.Status = EmailStatusPending
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[m.ID] = m
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*EmailMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id string, status EmailStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, ok := r.data[id]
	if !ok {
		return nil
	}
	m.Status = status
	m.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryRepository) ListByRecipient(ctx context.Context, email string, limit, offset int) ([]*EmailMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*EmailMessage
	for _, m := range r.data {
		for _, to := range m.To {
			if to == email {
				result = append(result, m)
				break
			}
		}
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	if offset > len(result) {
		return []*EmailMessage{}, nil
	}

	return result[offset:end], nil
}
