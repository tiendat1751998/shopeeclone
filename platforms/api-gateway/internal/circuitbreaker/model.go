package circuitbreaker

import "time"

type State string

const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateHalfOpen State = "half_open"
)

type CircuitBreaker struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	ServiceName          string    `json:"service_name"`
	FailureThreshold     int       `json:"failure_threshold"`
	RecoveryTimeout      int       `json:"recovery_timeout_seconds"`
	HalfOpenMaxRequests  int       `json:"half_open_max_requests"`
	State                State     `json:"state"`
	FailureCount         int       `json:"failure_count"`
	LastFailure          time.Time `json:"last_failure"`
	HalfOpenSuccessCount int       `json:"half_open_success_count"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
