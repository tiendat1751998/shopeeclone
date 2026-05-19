package cohort

import (
	"context"
	"sync"
	"time"
)

type Repository interface {
	StoreAnalysis(ctx context.Context, analysis *CohortAnalysis) error
	GetAnalysis(ctx context.Context, id string) (*CohortAnalysis, error)
	ListAnalyses(ctx context.Context, offset, limit int) ([]*CohortAnalysis, int, error)
	StoreDefinition(ctx context.Context, def *CohortDefinition) error
	GetDefinition(ctx context.Context, id string) (*CohortDefinition, error)
}

type InMemoryRepository struct {
	mu          sync.RWMutex
	analyses    map[string]*CohortAnalysis
	definitions map[string]*CohortDefinition
	analysesList []*CohortAnalysis
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		analyses:     make(map[string]*CohortAnalysis),
		definitions:  make(map[string]*CohortDefinition),
		analysesList: make([]*CohortAnalysis, 0),
	}
}

func (r *InMemoryRepository) StoreAnalysis(ctx context.Context, analysis *CohortAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if analysis.AnalyzedAt.IsZero() {
		analysis.AnalyzedAt = time.Now()
	}
	r.analyses[analysis.ID] = analysis
	r.analysesList = append(r.analysesList, analysis)
	return nil
}

func (r *InMemoryRepository) GetAnalysis(ctx context.Context, id string) (*CohortAnalysis, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	analysis, ok := r.analyses[id]
	if !ok {
		return nil, nil
	}
	return analysis, nil
}

func (r *InMemoryRepository) ListAnalyses(ctx context.Context, offset, limit int) ([]*CohortAnalysis, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.analysesList)
	if offset >= total {
		return []*CohortAnalysis{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return r.analysesList[offset:end], total, nil
}

func (r *InMemoryRepository) StoreDefinition(ctx context.Context, def *CohortDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if def.CreatedAt.IsZero() {
		def.CreatedAt = time.Now()
	}
	r.definitions[def.ID] = def
	return nil
}

func (r *InMemoryRepository) GetDefinition(ctx context.Context, id string) (*CohortDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	def, ok := r.definitions[id]
	if !ok {
		return nil, nil
	}
	return def, nil
}
