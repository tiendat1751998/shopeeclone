package slo

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Window string

const (
	Window1d  Window = "1d"
	Window7d  Window = "7d"
	Window28d Window = "28d"
)

type SLO struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Service          string  `json:"service"`
	SLIMetric        string  `json:"sli_metric"`
	TargetPercentage float64 `json:"target_percentage"`
	Window           Window  `json:"window"`
	CurrentValue     float64 `json:"current_value"`
	BudgetRemaining  float64 `json:"budget_remaining"`
}

type SLI struct {
	MetricName string  `json:"metric_name"`
	GoodCount  int64   `json:"good_count"`
	TotalCount int64   `json:"total_count"`
	Ratio      float64 `json:"ratio"`
}

type SLOReport struct {
	SLO            *SLO   `json:"slo"`
	SLI            *SLI   `json:"sli"`
	Compliant      bool   `json:"compliant"`
	BudgetUsed     float64 `json:"budget_used"`
	ErrorBudget    float64 `json:"error_budget"`
}

type Repository interface {
	Create(slo *SLO) error
	Get(id string) (*SLO, error)
	List() ([]*SLO, error)
	Update(slo *SLO) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*SLO
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*SLO),
	}
}

var ErrSLONotFound = errors.New("SLO not found")

func (r *InMemoryRepository) Create(slo *SLO) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	slo.ID = uuid.New().String()
	r.items[slo.ID] = slo
	return nil
}

func (r *InMemoryRepository) Get(id string) (*SLO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	slo, ok := r.items[id]
	if !ok {
		return nil, ErrSLONotFound
	}
	return slo, nil
}

func (r *InMemoryRepository) List() ([]*SLO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*SLO, 0, len(r.items))
	for _, slo := range r.items {
		result = append(result, slo)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(slo *SLO) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[slo.ID]; !ok {
		return ErrSLONotFound
	}
	r.items[slo.ID] = slo
	return nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSLO(name, service, sliMetric string, targetPercentage float64, window Window) (*SLO, error) {
	slo := &SLO{
		Name:             name,
		Service:          service,
		SLIMetric:        sliMetric,
		TargetPercentage: targetPercentage,
		Window:           window,
		CurrentValue:     100.0,
		BudgetRemaining:  100.0,
	}
	if err := s.repo.Create(slo); err != nil {
		return nil, err
	}
	return slo, nil
}

func (s *Service) CalculateSLI(metricName string, goodCount, totalCount int64) *SLI {
	var ratio float64
	if totalCount > 0 {
		ratio = float64(goodCount) / float64(totalCount) * 100
	}
	return &SLI{
		MetricName: metricName,
		GoodCount:  goodCount,
		TotalCount: totalCount,
		Ratio:      ratio,
	}
}

func (s *Service) EvaluateSLO(sloID string, sli *SLI) (*SLO, error) {
	slo, err := s.repo.Get(sloID)
	if err != nil {
		return nil, err
	}
	slo.CurrentValue = sli.Ratio
	budget := slo.TargetPercentage / 100.0
	if sli.Ratio/100.0 > 0 {
		slo.BudgetRemaining = ((sli.Ratio / 100.0) - (1 - budget)) / budget * 100
	} else {
		slo.BudgetRemaining = 0
	}
	if slo.BudgetRemaining < 0 {
		slo.BudgetRemaining = 0
	}
	if err := s.repo.Update(slo); err != nil {
		return nil, err
	}
	return slo, nil
}

func (s *Service) GetSLOReport(sloID string) (*SLOReport, error) {
	slo, err := s.repo.Get(sloID)
	if err != nil {
		return nil, err
	}
	sli := &SLI{
		MetricName: slo.SLIMetric,
		Ratio:      slo.CurrentValue,
	}
	report := &SLOReport{
		SLO: slo,
		SLI: sli,
	}
	report.Compliant = slo.CurrentValue >= slo.TargetPercentage
	errorBudget := 100.0 - slo.TargetPercentage
	report.ErrorBudget = errorBudget
	if slo.CurrentValue > 0 {
		used := (100.0 - slo.CurrentValue)
		report.BudgetUsed = (used / errorBudget) * 100
	} else {
		report.BudgetUsed = 0
	}
	return report, nil
}

func (s *Service) ListSLOs() ([]*SLO, error) {
	return s.repo.List()
}
