package itemembedding

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	Store(ctx context.Context, emb *ItemEmbedding) error
	Get(ctx context.Context, itemID string) (*ItemEmbedding, error)
	GetAll(ctx context.Context) ([]*ItemEmbedding, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	entries map[string]*ItemEmbedding
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		entries: make(map[string]*ItemEmbedding),
	}
}

func (r *InMemoryRepository) Store(ctx context.Context, emb *ItemEmbedding) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	emb.UpdatedAt = time.Now()
	r.entries[emb.ItemID] = emb
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, itemID string) (*ItemEmbedding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	emb, ok := r.entries[itemID]
	if !ok {
		return nil, ErrItemEmbeddingNotFound
	}
	return emb, nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]*ItemEmbedding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*ItemEmbedding, 0, len(r.entries))
	for _, emb := range r.entries {
		result = append(result, emb)
	}
	return result, nil
}
