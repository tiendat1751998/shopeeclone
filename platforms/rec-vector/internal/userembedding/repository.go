package userembedding

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Store(ctx context.Context, emb *UserEmbedding) error
	Get(ctx context.Context, userID string) (*UserEmbedding, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	entries map[string]*UserEmbedding
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		entries: make(map[string]*UserEmbedding),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, emb *UserEmbedding) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	emb.UpdatedAt = time.Now()
	r.entries[emb.UserID] = emb
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, userID string) (*UserEmbedding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	emb, ok := r.entries[userID]
	if !ok {
		return nil, ErrUserEmbeddingNotFound
	}
	return emb, nil
}
