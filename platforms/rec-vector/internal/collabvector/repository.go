package collabvector

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Repository interface {
	StoreInteraction(ctx context.Context, interaction *Interaction) error
	GetInteractions(ctx context.Context) ([]*Interaction, error)
	GetMatrix(ctx context.Context) (*InteractionMatrix, error)
	StoreFactors(ctx context.Context, userFactors, itemFactors [][]float64) error
	GetFactors(ctx context.Context) (*LatentFactors, error)
}

type InMemoryRepository struct {
	mu           sync.RWMutex
	interactions []*Interaction
	userFactors  [][]float64
	itemFactors  [][]float64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		interactions: make([]*Interaction, 0),
	}
}

func (r *InMemoryRepository) StoreInteraction(ctx context.Context, interaction *Interaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if interaction.Timestamp.IsZero() {
		interaction.Timestamp = time.Now()
	}
	r.interactions = append(r.interactions, interaction)
	return nil
}

func (r *InMemoryRepository) GetInteractions(ctx context.Context) ([]*Interaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Interaction, len(r.interactions))
	copy(result, r.interactions)
	return result, nil
}

func (r *InMemoryRepository) GetMatrix(ctx context.Context) (*InteractionMatrix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userItem := make(map[string]map[string]float64)
	userSet := make(map[string]bool)
	itemSet := make(map[string]bool)

	for _, inter := range r.interactions {
		if userItem[inter.UserID] == nil {
			userItem[inter.UserID] = make(map[string]float64)
		}
		userItem[inter.UserID][inter.ItemID] += inter.Weight
		userSet[inter.UserID] = true
		itemSet[inter.ItemID] = true
	}

	users := make([]string, 0, len(userSet))
	for u := range userSet {
		users = append(users, u)
	}
	items := make([]string, 0, len(itemSet))
	for i := range itemSet {
		items = append(items, i)
	}

	userIndex := make(map[string]int)
	for i, u := range users {
		userIndex[u] = i
	}
	itemIndex := make(map[string]int)
	for i, it := range items {
		itemIndex[it] = i
	}

	return &InteractionMatrix{
		UserItemRatings: userItem,
		UserIndex:       userIndex,
		ItemIndex:       itemIndex,
		Users:           users,
		Items:           items,
	}, nil
}

func (r *InMemoryRepository) StoreFactors(ctx context.Context, userFactors, itemFactors [][]float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userFactors = userFactors
	r.itemFactors = itemFactors
	return nil
}

func (r *InMemoryRepository) GetFactors(ctx context.Context) (*LatentFactors, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.userFactors) == 0 || len(r.itemFactors) == 0 {
		return nil, ErrNotEnoughInteractions
	}
	uf := make([][]float64, len(r.userFactors))
	for i := range r.userFactors {
		uf[i] = make([]float64, len(r.userFactors[i]))
		copy(uf[i], r.userFactors[i])
	}
	itf := make([][]float64, len(r.itemFactors))
	for i := range r.itemFactors {
		itf[i] = make([]float64, len(r.itemFactors[i]))
		copy(itf[i], r.itemFactors[i])
	}
	return &LatentFactors{UserFactors: uf, ItemFactors: itf}, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
