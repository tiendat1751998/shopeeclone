package circuitbreaker

import (
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, serviceName string, failureThreshold, recoveryTimeout, halfOpenMaxRequests int) (*CircuitBreaker, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if serviceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if failureThreshold <= 0 {
		failureThreshold = 5
	}
	if recoveryTimeout < 0 {
		recoveryTimeout = 60
	}
	if halfOpenMaxRequests <= 0 {
		halfOpenMaxRequests = 3
	}

	now := time.Now()
	cb := &CircuitBreaker{
		ID:                  fmt.Sprintf("cb-%s", name),
		Name:                name,
		ServiceName:         serviceName,
		FailureThreshold:    failureThreshold,
		RecoveryTimeout:     recoveryTimeout,
		HalfOpenMaxRequests: halfOpenMaxRequests,
		State:               StateClosed,
		FailureCount:        0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.Store(cb); err != nil {
		return nil, err
	}
	return cb, nil
}

func (s *Service) RecordSuccess(id string) error {
	cb, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if cb == nil {
		return fmt.Errorf("circuit breaker not found: %s", id)
	}

	switch cb.State {
	case StateHalfOpen:
		cb.HalfOpenSuccessCount++
		if cb.HalfOpenSuccessCount >= cb.HalfOpenMaxRequests {
			cb.State = StateClosed
			cb.FailureCount = 0
			cb.HalfOpenSuccessCount = 0
		}
	case StateClosed:
		cb.FailureCount = 0
	}

	cb.UpdatedAt = time.Now()
	return s.repo.Store(cb)
}

func (s *Service) RecordFailure(id string) error {
	cb, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if cb == nil {
		return fmt.Errorf("circuit breaker not found: %s", id)
	}

	now := time.Now()
	cb.FailureCount++
	cb.LastFailure = now

	switch cb.State {
	case StateClosed:
		if cb.FailureCount >= cb.FailureThreshold {
			cb.State = StateOpen
			cb.HalfOpenSuccessCount = 0
		}
	case StateHalfOpen:
		cb.State = StateOpen
		cb.HalfOpenSuccessCount = 0
	}

	cb.UpdatedAt = now
	return s.repo.Store(cb)
}

func (s *Service) CanPass(id string) (bool, error) {
	cb, err := s.repo.Get(id)
	if err != nil {
		return false, err
	}
	if cb == nil {
		return true, nil
	}

	switch cb.State {
	case StateClosed:
		return true, nil
	case StateOpen:
		elapsed := time.Since(cb.LastFailure).Seconds()
		if elapsed >= float64(cb.RecoveryTimeout) {
			cb.State = StateHalfOpen
			cb.HalfOpenSuccessCount = 0
			cb.UpdatedAt = time.Now()
			s.repo.Store(cb)
			return true, nil
		}
		return false, nil
	case StateHalfOpen:
		return cb.HalfOpenSuccessCount < cb.HalfOpenMaxRequests, nil
	}

	return true, nil
}

func (s *Service) GetState(id string) (State, error) {
	cb, err := s.repo.Get(id)
	if err != nil {
		return "", err
	}
	if cb == nil {
		return StateClosed, nil
	}
	return cb.State, nil
}

func (s *Service) List() ([]*CircuitBreaker, error) {
	return s.repo.List()
}

func (s *Service) Get(id string) (*CircuitBreaker, error) {
	return s.repo.Get(id)
}
