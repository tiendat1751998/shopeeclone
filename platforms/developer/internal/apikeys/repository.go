package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type Repository interface {
	Store(ctx context.Context, key *APIKey) error
	GetByKeyHash(ctx context.Context, hash string) (*APIKey, error)
	GetByID(ctx context.Context, id string) (*APIKey, error)
	List(ctx context.Context) ([]*APIKey, error)
	Update(ctx context.Context, key *APIKey) error
	DeleteByHash(ctx context.Context, hash string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	byID  map[string]*APIKey
	byKey map[string]string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		byID:  make(map[string]*APIKey),
		byKey: make(map[string]string),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, key *APIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[key.ID] = key
	r.byKey[key.KeyHash] = key.ID
	return nil
}

func (r *InMemoryRepository) GetByKeyHash(ctx context.Context, hash string) (*APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byKey[hash]
	if !ok {
		return nil, nil
	}
	return r.byID[id], nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return key, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]*APIKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*APIKey
	for _, k := range r.byID {
		result = append(result, k)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, key *APIKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[key.ID] = key
	r.byKey[key.KeyHash] = key.ID
	return nil
}

func (r *InMemoryRepository) DeleteByHash(ctx context.Context, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byKey, hash)
	return nil
}

func HashKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(h[:])
}
