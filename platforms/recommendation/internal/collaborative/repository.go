package collaborative

import (
	"context"
	"sync"
)

type Repository interface {
	GetUserRatings(ctx context.Context, userID string) (map[string]float64, error)
	GetItemRatings(ctx context.Context, itemID string) (map[string]float64, error)
	GetAllUserIDs(ctx context.Context) ([]string, error)
	GetAllItemIDs(ctx context.Context) ([]string, error)
	StoreRating(ctx context.Context, userID, itemID string, rating float64) error
	GetMatrix(ctx context.Context) (*RatingMatrix, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	userRatings map[string]map[string]float64
	itemRatings map[string]map[string]float64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		userRatings: make(map[string]map[string]float64),
		itemRatings: make(map[string]map[string]float64),
	}
}

func (r *InMemoryRepository) GetUserRatings(ctx context.Context, userID string) (map[string]float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ratings, ok := r.userRatings[userID]
	if !ok {
		return map[string]float64{}, nil
	}
	result := make(map[string]float64)
	for k, v := range ratings {
		result[k] = v
	}
	return result, nil
}

func (r *InMemoryRepository) GetItemRatings(ctx context.Context, itemID string) (map[string]float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ratings, ok := r.itemRatings[itemID]
	if !ok {
		return map[string]float64{}, nil
	}
	result := make(map[string]float64)
	for k, v := range ratings {
		result[k] = v
	}
	return result, nil
}

func (r *InMemoryRepository) GetAllUserIDs(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.userRatings))
	for id := range r.userRatings {
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *InMemoryRepository) GetAllItemIDs(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.itemRatings))
	for id := range r.itemRatings {
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *InMemoryRepository) StoreRating(ctx context.Context, userID, itemID string, rating float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.userRatings[userID] == nil {
		r.userRatings[userID] = make(map[string]float64)
	}
	r.userRatings[userID][itemID] = rating
	if r.itemRatings[itemID] == nil {
		r.itemRatings[itemID] = make(map[string]float64)
	}
	r.itemRatings[itemID][userID] = rating
	return nil
}

func (r *InMemoryRepository) GetMatrix(ctx context.Context) (*RatingMatrix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	userCopy := make(map[string]map[string]float64)
	for u, items := range r.userRatings {
		userCopy[u] = make(map[string]float64)
		for i, v := range items {
			userCopy[u][i] = v
		}
	}
	itemCopy := make(map[string]map[string]float64)
	for i, users := range r.itemRatings {
		itemCopy[i] = make(map[string]float64)
		for u, v := range users {
			itemCopy[i][u] = v
		}
	}
	return &RatingMatrix{UserItemRatings: userCopy, ItemUserRatings: itemCopy}, nil
}
