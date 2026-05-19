package scoring

import "sync"

type Repository interface {
	GetThresholds() *ScoreThreshold
	SetThresholds(t *ScoreThreshold)
}

type InMemoryRepository struct {
	mu         sync.RWMutex
	thresholds *ScoreThreshold
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		thresholds: &ScoreThreshold{
			Low:      0,
			Medium:   26,
			High:     51,
			Critical: 76,
		},
	}
}

func (r *InMemoryRepository) GetThresholds() *ScoreThreshold {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.thresholds
}

func (r *InMemoryRepository) SetThresholds(t *ScoreThreshold) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.thresholds = t
}
