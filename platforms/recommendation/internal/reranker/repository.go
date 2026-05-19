package reranker

import (
	"context"
	"sync"
)

type ExposureRecord struct {
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}

type Repository interface {
	GetExposureCount(ctx context.Context, productID string) (int, error)
	IncrementExposure(ctx context.Context, productID string) error
	GetConfig(ctx context.Context) (*ReRankConfig, error)
	SetConfig(ctx context.Context, cfg *ReRankConfig) error
}

type ReRankConfig struct {
	MaxPerCategory      int     `json:"max_per_category"`
	NewItemBoost        float64 `json:"new_item_boost"`
	NewItemHours        int     `json:"new_item_hours"`
	ExposureDownrank    float64 `json:"exposure_downrank"`
	DiversityWeight     float64 `json:"diversity_weight"`
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	exposure map[string]int
	config   *ReRankConfig
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		exposure: make(map[string]int),
		config:   DefaultReRankConfig(),
	}
}

func (r *InMemoryRepository) GetExposureCount(ctx context.Context, productID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.exposure[productID], nil
}

func (r *InMemoryRepository) IncrementExposure(ctx context.Context, productID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exposure[productID]++
	return nil
}

func (r *InMemoryRepository) GetConfig(ctx context.Context) (*ReRankConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg := *r.config
	return &cfg, nil
}

func (r *InMemoryRepository) SetConfig(ctx context.Context, cfg *ReRankConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cfgCopy := *cfg
	r.config = &cfgCopy
	return nil
}

func DefaultReRankConfig() *ReRankConfig {
	return &ReRankConfig{
		MaxPerCategory:      3,
		NewItemBoost:        0.2,
		NewItemHours:        168,
		ExposureDownrank:    0.05,
		DiversityWeight:     0.3,
	}
}
