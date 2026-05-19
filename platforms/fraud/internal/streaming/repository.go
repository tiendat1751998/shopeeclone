package streaming

import (
	"context"
	"sync"
	"time"
)

type windowEntry struct {
	timestamp time.Time
	value     float64
}

type Repository interface {
	AddEvent(ctx context.Context, entityType, entityID string, timestamp time.Time)
	CountSince(ctx context.Context, entityType, entityID string, since time.Time) int
	GetRecentEvents(ctx context.Context, entityType, entityID string, since time.Time) []windowEntry
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	events map[string][]windowEntry
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events: make(map[string][]windowEntry),
	}
}

func (r *InMemoryRepository) key(entityType, entityID string) string {
	return entityType + ":" + entityID
}

func (r *InMemoryRepository) AddEvent(ctx context.Context, entityType, entityID string, timestamp time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := r.key(entityType, entityID)
	r.events[k] = append(r.events[k], windowEntry{timestamp: timestamp})
}

func (r *InMemoryRepository) CountSince(ctx context.Context, entityType, entityID string, since time.Time) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k := r.key(entityType, entityID)
	entries, ok := r.events[k]
	if !ok {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.timestamp.After(since) {
			count++
		}
	}
	return count
}

func (r *InMemoryRepository) GetRecentEvents(ctx context.Context, entityType, entityID string, since time.Time) []windowEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k := r.key(entityType, entityID)
	entries, ok := r.events[k]
	if !ok {
		return nil
	}
	var result []windowEntry
	for _, e := range entries {
		if e.timestamp.After(since) {
			result = append(result, e)
		}
	}
	return result
}
