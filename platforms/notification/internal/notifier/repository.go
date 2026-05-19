package notifier

import (
	"context"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]*Notification, error)
	UpdateStatus(ctx context.Context, id string, status DeliveryStatus) error
	MarkRead(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

type InMemoryRepository struct {
	mu            sync.RWMutex
	notifications map[string]*Notification
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		notifications: make(map[string]*Notification),
	}
}

func (r *InMemoryRepository) Create(ctx context.Context, n *Notification) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	if n.Status == "" {
		n.Status = StatusPending
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[n.ID] = n
	return nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	n, ok := r.notifications[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (r *InMemoryRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*Notification
	for _, n := range r.notifications {
		if n.UserID == userID {
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
		return []*Notification{}, nil
	}

	return results[start:end], nil
}

func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id string, status DeliveryStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notifications[id]
	if !ok {
		return ErrNotFound
	}
	n.Status = status
	n.UpdatedAt = time.Now()
	return nil
}

func (r *InMemoryRepository) MarkRead(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notifications[id]
	if !ok {
		return ErrNotFound
	}
	now := time.Now()
	n.Status = StatusRead
	n.ReadAt = &now
	n.UpdatedAt = now
	return nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notifications[id]; !ok {
		return ErrNotFound
	}
	delete(r.notifications, id)
	return nil
}

func (r *InMemoryRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, n := range r.notifications {
		if n.UserID == userID && n.Status != StatusRead {
			count++
		}
	}
	return count, nil
}

func matchesCategoryType(doc *Notification, q SearchQuery) bool {
	if q.Type != "" && doc.Type != q.Type {
		return false
	}
	if q.Channel != "" && doc.Channel != q.Channel {
		return false
	}
	if q.Status != "" && doc.Status != q.Status {
		return false
	}
	if q.Query != "" {
		qLower := strings.ToLower(q.Query)
		titleLower := strings.ToLower(doc.Title)
		bodyLower := strings.ToLower(doc.Body)
		if !strings.Contains(titleLower, qLower) && !strings.Contains(bodyLower, qLower) {
			return false
		}
	}
	return true
}

type SearchQuery struct {
	Query   string
	Type    NotificationType
	Channel Channel
	Status  DeliveryStatus
	Page    int
	Limit   int
}
