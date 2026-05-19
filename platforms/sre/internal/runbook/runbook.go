package runbook

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Step struct {
	Title          string `json:"title"`
	Command        string `json:"command"`
	ExpectedResult string `json:"expected_result"`
}

type Runbook struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Service      string    `json:"service"`
	IncidentType string    `json:"incident_type"`
	Steps        []Step    `json:"steps"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Filter struct {
	Service      string `form:"service"`
	IncidentType string `form:"incident_type"`
}

type Repository interface {
	Create(rb *Runbook) error
	Get(id string) (*Runbook, error)
	Update(rb *Runbook) error
	Delete(id string) error
	List(filter Filter) ([]*Runbook, error)
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*Runbook
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*Runbook),
	}
}

var (
	ErrRunbookNotFound = errors.New("runbook not found")
)

func (r *InMemoryRepository) Create(rb *Runbook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rb.ID = uuid.New().String()
	rb.CreatedAt = time.Now()
	rb.UpdatedAt = time.Now()
	r.items[rb.ID] = rb
	return nil
}

func (r *InMemoryRepository) Get(id string) (*Runbook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rb, ok := r.items[id]
	if !ok {
		return nil, ErrRunbookNotFound
	}
	return rb, nil
}

func (r *InMemoryRepository) Update(rb *Runbook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[rb.ID]; !ok {
		return ErrRunbookNotFound
	}
	rb.UpdatedAt = time.Now()
	r.items[rb.ID] = rb
	return nil
}

func (r *InMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrRunbookNotFound
	}
	delete(r.items, id)
	return nil
}

func (r *InMemoryRepository) List(filter Filter) ([]*Runbook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Runbook
	for _, rb := range r.items {
		if filter.Service != "" && rb.Service != filter.Service {
			continue
		}
		if filter.IncidentType != "" && rb.IncidentType != filter.IncidentType {
			continue
		}
		result = append(result, rb)
	}
	return result, nil
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(title, service, incidentType string, steps []Step) (*Runbook, error) {
	rb := &Runbook{
		Title:        title,
		Service:      service,
		IncidentType: incidentType,
		Steps:        steps,
	}
	if err := s.repo.Create(rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Service) Get(id string) (*Runbook, error) {
	return s.repo.Get(id)
}

func (s *Service) Update(id, title string, steps []Step) (*Runbook, error) {
	rb, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	rb.Title = title
	rb.Steps = steps
	if err := s.repo.Update(rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) List(filter Filter) ([]*Runbook, error) {
	return s.repo.List(filter)
}
