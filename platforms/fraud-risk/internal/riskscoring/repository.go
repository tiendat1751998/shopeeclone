package riskscoring

import (
	"sync"
)

type Repository interface {
	GetThresholds() *RiskLevelThresholds
	SetThresholds(t *RiskLevelThresholds)
	UpsertFactor(f RiskFactor)
	ListFactors() []RiskFactor
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	thresholds *RiskLevelThresholds
	factors    map[string]RiskFactor
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		thresholds: &RiskLevelThresholds{
			Safe:     20,
			Low:      40,
			Medium:   60,
			High:     80,
			Critical: 100,
		},
		factors: make(map[string]RiskFactor),
	}
}

func (r *InMemoryRepository) GetThresholds() *RiskLevelThresholds {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.thresholds
}

func (r *InMemoryRepository) SetThresholds(t *RiskLevelThresholds) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.thresholds = t
}

func (r *InMemoryRepository) UpsertFactor(f RiskFactor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factors[f.Name] = f
}

func (r *InMemoryRepository) ListFactors() []RiskFactor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]RiskFactor, 0, len(r.factors))
	for _, f := range r.factors {
		result = append(result, f)
	}
	return result
}
