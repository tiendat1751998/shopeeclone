package blacklist

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, entry *BlacklistEntry) error
	Remove(ctx context.Context, id string) error
	GetByTypeAndValue(ctx context.Context, bt BlacklistType, value string) (*BlacklistEntry, error)
	List(ctx context.Context) ([]*BlacklistEntry, error)
	ListAll(ctx context.Context) ([]*BlacklistEntry, error)
	Update(ctx context.Context, entry *BlacklistEntry) error
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	entries map[string]*BlacklistEntry
	byValue map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		entries: make(map[string]*BlacklistEntry),
		byValue: make(map[string]string),
	}
}

func (r *InMemoryRepository) key(bt BlacklistType, value string) string {
	return string(bt) + ":" + value
}

func (r *InMemoryRepository) Add(ctx context.Context, entry *BlacklistEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	entry.IsActive = true
	r.entries[entry.ID] = entry
	r.byValue[r.key(entry.Type, entry.Value)] = entry.ID
	return nil
}

func (r *InMemoryRepository) Remove(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[id]
	if !ok {
		return ErrEntryNotFound
	}
	delete(r.byValue, r.key(entry.Type, entry.Value))
	delete(r.entries, id)
	return nil
}

func (r *InMemoryRepository) GetByTypeAndValue(ctx context.Context, bt BlacklistType, value string) (*BlacklistEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byValue[r.key(bt, value)]
	if !ok {
		return nil, ErrEntryNotFound
	}
	entry, ok := r.entries[id]
	if !ok {
		return nil, ErrEntryNotFound
	}
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return nil, ErrEntryNotFound
	}
	return entry, nil
}

func (r *InMemoryRepository) ListAll(ctx context.Context) ([]*BlacklistEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*BlacklistEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		result = append(result, entry)
	}
	return result, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*BlacklistEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*BlacklistEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, entry *BlacklistEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
	return nil
}
