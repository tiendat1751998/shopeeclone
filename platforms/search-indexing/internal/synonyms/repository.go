package synonyms

import (
	"context"
	"sync"
)

type Repository interface {
	CreateSet(ctx context.Context, set *SynonymSet) error
	GetSet(ctx context.Context, id string) (*SynonymSet, error)
	ListSets(ctx context.Context) ([]*SynonymSet, error)
	UpdateSet(ctx context.Context, set *SynonymSet) error
	DeleteSet(ctx context.Context, id string) error
	GetGraph(ctx context.Context) (*SynonymGraph, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	sets  map[string]*SynonymSet
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		sets: make(map[string]*SynonymSet),
	}
}

func (r *InMemoryRepository) CreateSet(_ context.Context, set *SynonymSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.sets[set.ID]; ok {
		return ErrSetAlreadyExists
	}
	r.sets[set.ID] = set
	return nil
}

func (r *InMemoryRepository) GetSet(_ context.Context, id string) (*SynonymSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	set, ok := r.sets[id]
	if !ok {
		return nil, ErrSynonymSetNotFound
	}
	return set, nil
}

func (r *InMemoryRepository) ListSets(_ context.Context) ([]*SynonymSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*SynonymSet, 0, len(r.sets))
	for _, s := range r.sets {
		list = append(list, s)
	}
	return list, nil
}

func (r *InMemoryRepository) UpdateSet(_ context.Context, set *SynonymSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.sets[set.ID]; !ok {
		return ErrSynonymSetNotFound
	}
	r.sets[set.ID] = set
	return nil
}

func (r *InMemoryRepository) DeleteSet(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.sets[id]; !ok {
		return ErrSynonymSetNotFound
	}
	delete(r.sets, id)
	return nil
}

func (r *InMemoryRepository) GetGraph(_ context.Context) (*SynonymGraph, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	edges := make(map[string][]string)
	for _, set := range r.sets {
		if !set.IsActive {
			continue
		}
		for _, word := range set.Words {
			for _, other := range set.Words {
				if word != other {
					edges[word] = append(edges[word], other)
				}
			}
		}
	}
	return &SynonymGraph{Edges: edges}, nil
}
