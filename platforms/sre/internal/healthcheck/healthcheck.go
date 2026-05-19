package healthcheck

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPass    Status = "pass"
	StatusFail    Status = "fail"
	StatusDegraded Status = "degraded"
)

type HealthCheck struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Target          string `json:"target"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	Method          string `json:"method"`
	ExpectedStatus  int    `json:"expected_status"`
}

type CheckResult struct {
	ID             string    `json:"id"`
	CheckName      string    `json:"check_name"`
	Status         Status    `json:"status"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	Error          string    `json:"error,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

type Repository interface {
	CreateCheck(hc *HealthCheck) error
	ListChecks() ([]*HealthCheck, error)
	AddResult(result *CheckResult) error
	GetResults() ([]*CheckResult, error)
}

type InMemoryRepository struct {
	mu      sync.RWMutex
	checks  map[string]*HealthCheck
	results []*CheckResult
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		checks:  make(map[string]*HealthCheck),
		results: make([]*CheckResult, 0),
	}
}

var ErrCheckNotFound = errors.New("health check not found")

func (r *InMemoryRepository) CreateCheck(hc *HealthCheck) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	hc.ID = uuid.New().String()
	r.checks[hc.ID] = hc
	return nil
}

func (r *InMemoryRepository) ListChecks() ([]*HealthCheck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*HealthCheck, 0, len(r.checks))
	for _, hc := range r.checks {
		result = append(result, hc)
	}
	return result, nil
}

func (r *InMemoryRepository) AddResult(result *CheckResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	result.ID = uuid.New().String()
	result.CheckedAt = time.Now()
	r.results = append(r.results, result)
	return nil
}

func (r *InMemoryRepository) GetResults() ([]*CheckResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*CheckResult, len(r.results))
	copy(result, r.results)
	return result, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCheck(name, target string, interval, timeout int, method string, expectedStatus int) (*HealthCheck, error) {
	hc := &HealthCheck{
		Name:            name,
		Target:          target,
		IntervalSeconds: interval,
		TimeoutSeconds:  timeout,
		Method:          method,
		ExpectedStatus:  expectedStatus,
	}
	if err := s.repo.CreateCheck(hc); err != nil {
		return nil, err
	}
	return hc, nil
}

func (s *Service) RunChecks() ([]*CheckResult, error) {
	checks, err := s.repo.ListChecks()
	if err != nil {
		return nil, err
	}
	var results []*CheckResult
	for _, hc := range checks {
		result := s.simulateCheck(hc)
		s.repo.AddResult(result)
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) simulateCheck(hc *HealthCheck) *CheckResult {
	r := rand.Intn(100)
	status := StatusPass
	rt := int64(rand.Intn(200) + 10)
	var errStr string
	if r < 10 {
		status = StatusFail
		rt = 0
		errStr = "connection timeout"
	} else if r < 20 {
		status = StatusDegraded
		rt = int64(rand.Intn(800) + 500)
		errStr = "slow response"
	}
	return &CheckResult{
		CheckName:      hc.Name,
		Status:         status,
		ResponseTimeMs: rt,
		Error:          errStr,
	}
}

func (s *Service) GetResults() ([]*CheckResult, error) {
	return s.repo.GetResults()
}

func (s *Service) ListChecks() ([]*HealthCheck, error) {
	return s.repo.ListChecks()
}
