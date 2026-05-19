package alerting

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AlertStatus string

const (
	AlertFiring       AlertStatus = "firing"
	AlertResolved     AlertStatus = "resolved"
	AlertAcknowledged AlertStatus = "acknowledged"
)

type AlertSeverity string

const (
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityInfo     AlertSeverity = "info"
)

type Alert struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Condition   string        `json:"condition"`
	Threshold   float64       `json:"threshold"`
	CurrentValue float64      `json:"current_value"`
	Status      AlertStatus   `json:"status"`
	Severity    AlertSeverity `json:"severity"`
	Service     string        `json:"service"`
	TriggeredAt time.Time     `json:"triggered_at"`
	ResolvedAt  *time.Time    `json:"resolved_at,omitempty"`
}

type Rule struct {
	Name            string  `json:"name"`
	MetricName      string  `json:"metric_name"`
	Operator        string  `json:"operator"`
	Threshold       float64 `json:"threshold"`
	DurationSeconds int     `json:"duration_seconds"`
	CooldownSeconds int     `json:"cooldown_seconds"`
}

type MetricValue struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type Repository interface {
	CreateAlert(a *Alert) error
	GetAlert(id string) (*Alert, error)
	ListAlerts() ([]*Alert, error)
	UpdateAlert(a *Alert) error
	CreateRule(r *Rule) error
	ListRules() ([]*Rule, error)
	GetRule(name string) (*Rule, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	alerts map[string]*Alert
	rules  map[string]*Rule
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		alerts: make(map[string]*Alert),
		rules:  make(map[string]*Rule),
	}
}

var (
	ErrAlertNotFound = errors.New("alert not found")
	ErrRuleNotFound  = errors.New("rule not found")
	ErrRuleExists    = errors.New("rule already exists")
)

func (r *InMemoryRepository) CreateAlert(a *Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a.ID = uuid.New().String()
	a.TriggeredAt = time.Now()
	r.alerts[a.ID] = a
	return nil
}

func (r *InMemoryRepository) GetAlert(id string) (*Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.alerts[id]
	if !ok {
		return nil, ErrAlertNotFound
	}
	return a, nil
}

func (r *InMemoryRepository) ListAlerts() ([]*Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Alert, 0, len(r.alerts))
	for _, a := range r.alerts {
		result = append(result, a)
	}
	return result, nil
}

func (r *InMemoryRepository) UpdateAlert(a *Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.alerts[a.ID]; !ok {
		return ErrAlertNotFound
	}
	r.alerts[a.ID] = a
	return nil
}

func (r *InMemoryRepository) CreateRule(rule *Rule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.rules[rule.Name]; ok {
		return ErrRuleExists
	}
	r.rules[rule.Name] = rule
	return nil
}

func (r *InMemoryRepository) ListRules() ([]*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result, nil
}

func (r *InMemoryRepository) GetRule(name string) (*Rule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.rules[name]
	if !ok {
		return nil, ErrRuleNotFound
	}
	return rule, nil
}
