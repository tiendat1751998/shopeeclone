package inapp

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, n *InAppNotification) error
	GetByID(ctx context.Context, id string) (*InAppNotification, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]*InAppNotification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Dismiss(ctx context.Context, id string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu            sync.RWMutex
	notifications map[string]*InAppNotification
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		notifications: make(map[string]*InAppNotification),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, n *InAppNotification) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	n.CreatedAt = time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[n.ID] = n
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*InAppNotification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	n, ok := r.notifications[id]
	if !ok {
		return nil, nil
	}
	return n, nil
}

func (r *InMemoryRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*InAppNotification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*InAppNotification
	for _, n := range r.notifications {
		if n.UserID == userID && !n.Dismissed {
			results = append(results, n)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	total := len(results)
	start := int(math.Min(float64(offset), float64(total)))
	end := int(math.Min(float64(offset+limit), float64(total)))

	if start >= total {
		return []*InAppNotification{}, nil
	}

	return results[start:end], nil
}

func (r *InMemoryRepository) MarkRead(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notifications[id]
	if !ok {
		return nil
	}
	now := time.Now()
	n.Read = true
	n.ReadAt = &now
	return nil
}

func (r *InMemoryRepository) MarkAllRead(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, n := range r.notifications {
		if n.UserID == userID {
			n.Read = true
			n.ReadAt = &now
		}
	}
	return nil
}

func (r *InMemoryRepository) Dismiss(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notifications[id]
	if !ok {
		return nil
	}
	n.Dismissed = true
	return nil
}

func (r *InMemoryRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, n := range r.notifications {
		if n.UserID == userID && !n.Read && !n.Dismissed {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.notifications, id)
	return nil
}
