package sms

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, m *SMSMessage) error
	GetByID(ctx context.Context, id string) (*SMSMessage, error)
	ListByRecipient(ctx context.Context, phone string, limit, offset int) ([]*SMSMessage, error)
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*SMSMessage
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*SMSMessage),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, m *SMSMessage) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.Status == "" {
		m.Status = "pending"
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[m.ID] = m
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*SMSMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (r *InMemoryRepository) ListByRecipient(ctx context.Context, phone string, limit, offset int) ([]*SMSMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*SMSMessage
	for _, m := range r.data {
		if m.To == phone {
			result = append(result, m)
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
		return []*SMSMessage{}, nil
	}

	return result[offset:end], nil
}
