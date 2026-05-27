package incident

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityMajor    Severity = "major"
	SeverityMinor    Severity = "minor"
)

type Status string

const (
	StatusDetected    Status = "detected"
	StatusTriaging    Status = "triaging"
	StatusMitigating  Status = "mitigating"
	StatusResolved    Status = "resolved"
	StatusPostmortem  Status = "postmortem"
)

type Incident struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Severity    Severity  `json:"severity"`
	Status      Status    `json:"status"`
	Service     string    `json:"service"`
	Region      string    `json:"region"`
	DetectedAt  time.Time `json:"detected_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Description string    `json:"description"`
	Assignee    string    `json:"assignee"`
}

type Filter struct {
	Status   Status   `form:"status"`
	Severity Severity `form:"severity"`
	Service  string   `form:"service"`
}

type Repository interface {
	Create(inc *Incident) error
	Get(id string) (*Incident, error)
	Update(inc *Incident) error
	List(filter Filter) ([]*Incident, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	items  map[string]*Incident
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		items: make(map[string]*Incident),
	}
}

var (
	ErrNotFound        = errors.New("incident not found")
	ErrInvalidSeverity = errors.New("invalid severity")
	ErrInvalidStatus   = errors.New("invalid status")
	ErrInvalidAssignee = errors.New("invalid assignee")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidTitle    = errors.New("invalid title")
)

func (r *InMemoryRepository) Create(inc *Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inc.ID = uuid.New().String()
	inc.Status = StatusDetected
	inc.DetectedAt = time.Now()
	r.items[inc.ID] = inc
	return nil
}

func (r *InMemoryRepository) Get(id string) (*Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inc, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return inc, nil
}

func (r *InMemoryRepository) Update(inc *Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[inc.ID]; !ok {
		return ErrNotFound
	}
	r.items[inc.ID] = inc
	return nil
}

func (r *InMemoryRepository) List(filter Filter) ([]*Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Incident
	for _, inc := range r.items {
		if filter.Status != "" && inc.Status != filter.Status {
			continue
		}
		if filter.Severity != "" && inc.Severity != filter.Severity {
			continue
		}
		if filter.Service != "" && inc.Service != filter.Service {
			continue
		}
		result = append(result, inc)
	}
	return result, nil
}
