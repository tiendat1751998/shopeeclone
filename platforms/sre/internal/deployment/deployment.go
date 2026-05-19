package deployment

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Strategy string

const (
	StrategyRolling  Strategy = "rolling"
	StrategyCanary   Strategy = "canary"
	StrategyBlueGreen Strategy = "bluegreen"
)

type DeployStatus string

const (
	DeployPending    DeployStatus = "pending"
	DeployInProgress DeployStatus = "in_progress"
	DeploySucceeded  DeployStatus = "succeeded"
	DeployFailed     DeployStatus = "failed"
	DeployRolledBack DeployStatus = "rolled_back"
)

type Deployment struct {
	ID                string       `json:"id"`
	Service           string       `json:"service"`
	Version           string       `json:"version"`
	Strategy          Strategy     `json:"strategy"`
	Status            DeployStatus `json:"status"`
	ProgressPercentage int         `json:"progress_percentage"`
	StartedAt         *time.Time   `json:"started_at,omitempty"`
	CompletedAt       *time.Time   `json:"completed_at,omitempty"`
}

var rollingSteps = []int{25, 50, 75, 100}
var canarySteps = []int{10, 25, 50, 100}

type Repository interface {
	Create(d *Deployment) error
	Get(id string) (*Deployment, error)
	Update(d *Deployment) error
	List() ([]*Deployment, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*Deployment
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*Deployment),
	}
}

var (
	ErrDeploymentNotFound = errors.New("deployment not found")
	ErrInvalidTransition  = errors.New("invalid status transition")
)

func (r *InMemoryRepository) Create(d *Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.ID = uuid.New().String()
	d.Status = DeployPending
	r.items[d.ID] = d
	return nil
}

func (r *InMemoryRepository) Get(id string) (*Deployment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.items[id]
	if !ok {
		return nil, ErrDeploymentNotFound
	}
	return d, nil
}

func (r *InMemoryRepository) Update(d *Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[d.ID]; !ok {
		return ErrDeploymentNotFound
	}
	r.items[d.ID] = d
	return nil
}

func (r *InMemoryRepository) List() ([]*Deployment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Deployment, 0, len(r.items))
	for _, d := range r.items {
		result = append(result, d)
	}
	return result, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(service, version string, strategy Strategy) (*Deployment, error) {
	now := time.Now()
	d := &Deployment{
		Service:  service,
		Version:  version,
		Strategy: strategy,
		Status:   DeployPending,
		StartedAt: &now,
	}
	if err := s.repo.Create(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Approve(id string) (*Deployment, error) {
	d, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if d.Status == DeploySucceeded || d.Status == DeployFailed || d.Status == DeployRolledBack {
		return nil, ErrInvalidTransition
	}
	d.Status = DeployInProgress
	var steps []int
	switch d.Strategy {
	case StrategyRolling:
		steps = rollingSteps
	case StrategyCanary:
		steps = canarySteps
	default:
		steps = []int{100}
	}
	for _, pct := range steps {
		if pct > d.ProgressPercentage {
			d.ProgressPercentage = pct
			break
		}
	}
	if d.ProgressPercentage == 100 {
		now := time.Now()
		d.Status = DeploySucceeded
		d.CompletedAt = &now
	}
	if err := s.repo.Update(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Rollback(id string) (*Deployment, error) {
	d, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if d.Status == DeployRolledBack {
		return nil, ErrInvalidTransition
	}
	d.Status = DeployRolledBack
	d.ProgressPercentage = 0
	now := time.Now()
	d.CompletedAt = &now
	if err := s.repo.Update(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetStatus(id string) (*Deployment, error) {
	return s.repo.Get(id)
}

func (s *Service) List() ([]*Deployment, error) {
	return s.repo.List()
}
