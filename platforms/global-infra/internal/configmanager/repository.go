package configmanager

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Create(ctx context.Context, entry *ConfigEntry) error
	Get(ctx context.Context, key string, env Environment, serviceName string) (*ConfigEntry, error)
	GetVersion(ctx context.Context, key string, env Environment, serviceName string, version int) (*ConfigEntry, error)
	List(ctx context.Context, serviceName string, env Environment) ([]*ConfigEntry, error)
	ListVersions(ctx context.Context, key string, env Environment, serviceName string) ([]*ConfigEntry, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	entries map[string]*ConfigEntry
	history map[string][]*ConfigEntry
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		entries: make(map[string]*ConfigEntry),
		history: make(map[string][]*ConfigEntry),
	}
}

func configKey(key string, env Environment, serviceName string) string {
	return serviceName + ":" + string(env) + ":" + key
}

func (r *InMemoryRepository) Create(ctx context.Context, entry *ConfigEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ck := configKey(entry.Key, entry.Environment, entry.ServiceName)
	entry.Version = 1
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	if existing, ok := r.entries[ck]; ok {
		entry.Version = existing.Version + 1
	}
	r.entries[ck] = entry
	hk := ck
	r.history[hk] = append(r.history[hk], entry)
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, key string, env Environment, serviceName string) (*ConfigEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[configKey(key, env, serviceName)]
	if !ok {
		return nil, nil
	}
	return entry, nil
}

func (r *InMemoryRepository) GetVersion(ctx context.Context, key string, env Environment, serviceName string, version int) (*ConfigEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hk := configKey(key, env, serviceName)
	versions := r.history[hk]
	for _, e := range versions {
		if e.Version == version {
			return e, nil
		}
	}
	return nil, nil
}

func (r *InMemoryRepository) List(ctx context.Context, serviceName string, env Environment) ([]*ConfigEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*ConfigEntry
	for _, e := range r.entries {
		if serviceName != "" && e.ServiceName != serviceName {
			continue
		}
		if env != "" && e.Environment != env {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *InMemoryRepository) ListVersions(ctx context.Context, key string, env Environment, serviceName string) ([]*ConfigEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hk := configKey(key, env, serviceName)
	versions := r.history[hk]
	return versions, nil
}
