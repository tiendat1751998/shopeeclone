package resilience

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker("test", CircuitBreakerOptions{
		MaxRequests:  5,
		Interval:     60 * time.Second,
		Timeout:      30 * time.Second,
		FailureRatio: 0.6,
		MinSamples:   5,
	})
	if cb == nil {
		t.Fatal("circuit breaker should not be nil")
	}
	if cb.Name() != "test" {
		t.Errorf("expected name 'test', got %s", cb.Name())
	}
	if cb.State() != StateClosed {
		t.Errorf("expected initial state closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_Defaults(t *testing.T) {
	cb := NewCircuitBreaker("defaults", CircuitBreakerOptions{})
	if cb.maxRequests != 5 {
		t.Errorf("expected default maxRequests 5, got %d", cb.maxRequests)
	}
	if cb.interval != 60*time.Second {
		t.Errorf("expected default interval 60s, got %v", cb.interval)
	}
	if cb.timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", cb.timeout)
	}
}

func TestCircuitBreaker_Success(t *testing.T) {
	cb := NewCircuitBreaker("test", CircuitBreakerOptions{MinSamples: 3})

	for i := 0; i < 5; i++ {
		err := cb.Execute(func() error {
			return nil
		})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	counts := cb.Counts()
	if counts.TotalSuccesses != 5 {
		t.Errorf("expected 5 successes, got %d", counts.TotalSuccesses)
	}
}

func TestCircuitBreaker_OpensOnFailures(t *testing.T) {
	cb := NewCircuitBreaker("test", CircuitBreakerOptions{
		MaxRequests:  1,
		Timeout:      50 * time.Millisecond,
		FailureRatio: 0.5,
		MinSamples:   3,
	})

	testErr := errors.New("test error")

	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return testErr
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("expected circuit breaker to be open, got %v", cb.State())
	}

	err := cb.Execute(func() error {
		return nil
	})
	if err == nil {
		t.Error("expected error when circuit breaker is open")
	}
}

func TestCircuitBreaker_HalfOpen(t *testing.T) {
	cb := NewCircuitBreaker("test", CircuitBreakerOptions{
		MaxRequests:  2,
		Timeout:      50 * time.Millisecond,
		FailureRatio: 0.5,
		MinSamples:   3,
	})

	testErr := errors.New("test error")

	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return testErr
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("expected open, got %v", cb.State())
	}

	time.Sleep(60 * time.Millisecond)

	if cb.State() != StateHalfOpen {
		t.Errorf("expected half-open after timeout, got %v", cb.State())
	}
}

func TestCircuitBreaker_ClosesAfterSuccesses(t *testing.T) {
	cb := NewCircuitBreaker("test", CircuitBreakerOptions{
		MaxRequests:  2,
		Timeout:      50 * time.Millisecond,
		FailureRatio: 0.5,
		MinSamples:   3,
	})

	testErr := errors.New("test error")
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return testErr
		})
	}

	time.Sleep(60 * time.Millisecond)

	if cb.State() != StateHalfOpen {
		t.Skip("expected half-open")
	}

	for i := 0; i < 2; i++ {
		cb.Execute(func() error {
			return nil
		})
	}

	if cb.State() != StateClosed {
		t.Errorf("expected closed after successful half-open probes, got %v", cb.State())
	}
}

func TestCircuitBreaker_Concurrency(t *testing.T) {
	cb := NewCircuitBreaker("concurrent", CircuitBreakerOptions{
		MinSamples: 100,
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb.Execute(func() error {
				return nil
			})
		}()
	}
	wg.Wait()

	counts := cb.Counts()
	if counts.TotalSuccesses != 50 {
		t.Errorf("expected 50 successes, got %d", counts.TotalSuccesses)
	}
}
