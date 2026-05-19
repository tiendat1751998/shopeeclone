package events

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreEvent(ctx context.Context, event *AnalyticsEvent) error
	GetEvent(ctx context.Context, eventID string) (*AnalyticsEvent, error)
	ListEvents(ctx context.Context, eventType EventType, startTime, endTime time.Time, offset, limit int) ([]*AnalyticsEvent, int64, error)
	ListEventsByUser(ctx context.Context, userID string, startTime, endTime time.Time, limit int) ([]*AnalyticsEvent, error)
	ListEventsBySession(ctx context.Context, sessionID string) ([]*AnalyticsEvent, error)
	GetEventCount(ctx context.Context, eventType EventType, startTime, endTime time.Time) (int64, error)
	GetRevenue(ctx context.Context, startTime, endTime time.Time) (float64, error)
	GetUniqueUsers(ctx context.Context, startTime, endTime time.Time) (int64, error)
	GetActiveUsers(ctx context.Context, startTime, endTime time.Time) (int64, error)
	GetOrders(ctx context.Context, startTime, endTime time.Time) (int64, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	events   map[string]*AnalyticsEvent
	eventIDs map[string]bool
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events:   make(map[string]*AnalyticsEvent),
		eventIDs: make(map[string]bool),
	}
}

func (r *InMemoryRepository) StoreEvent(ctx context.Context, event *AnalyticsEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if event.EventID == "" {
		return nil
	}
	if r.eventIDs[event.EventID] {
		return nil
	}
	r.eventIDs[event.EventID] = true
	r.events[event.EventID] = event
	return nil
}

func (r *InMemoryRepository) GetEvent(ctx context.Context, eventID string) (*AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	event, ok := r.events[eventID]
	if !ok {
		return nil, nil
	}
	return event, nil
}

func (r *InMemoryRepository) ListEvents(ctx context.Context, eventType EventType, startTime, endTime time.Time, offset, limit int) ([]*AnalyticsEvent, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*AnalyticsEvent
	for _, e := range r.events {
		if eventType != "" && e.EventType != eventType {
			continue
		}
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		filtered = append(filtered, e)
	}
	total := int64(len(filtered))
	if offset >= len(filtered) {
		return []*AnalyticsEvent{}, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryRepository) ListEventsByUser(ctx context.Context, userID string, startTime, endTime time.Time, limit int) ([]*AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*AnalyticsEvent
	for _, e := range r.events {
		if e.UserID != userID {
			continue
		}
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		filtered = append(filtered, e)
	}
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}

func (r *InMemoryRepository) ListEventsBySession(ctx context.Context, sessionID string) ([]*AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*AnalyticsEvent
	for _, e := range r.events {
		if e.SessionID == sessionID {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

func (r *InMemoryRepository) GetEventCount(ctx context.Context, eventType EventType, startTime, endTime time.Time) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, e := range r.events {
		if eventType != "" && e.EventType != eventType {
			continue
		}
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		count++
	}
	return count, nil
}

func (r *InMemoryRepository) GetRevenue(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var revenue float64
	for _, e := range r.events {
		if e.EventType != EventPurchase && e.EventType != EventCheckout {
			continue
		}
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		revenue += e.Revenue
	}
	return revenue, nil
}

func (r *InMemoryRepository) GetUniqueUsers(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make(map[string]bool)
	for _, e := range r.events {
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		users[e.UserID] = true
	}
	return int64(len(users)), nil
}

func (r *InMemoryRepository) GetActiveUsers(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make(map[string]bool)
	for _, e := range r.events {
		if e.EventType == EventPageview || e.EventType == EventClick || e.EventType == EventSearch || e.EventType == EventPurchase {
			if !startTime.IsZero() && e.Timestamp.Before(startTime) {
				continue
			}
			if !endTime.IsZero() && e.Timestamp.After(endTime) {
				continue
			}
			users[e.UserID] = true
		}
	}
	return int64(len(users)), nil
}

func (r *InMemoryRepository) GetOrders(ctx context.Context, startTime, endTime time.Time) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, e := range r.events {
		if e.EventType != EventPurchase {
			continue
		}
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		count++
	}
	return count, nil
}
