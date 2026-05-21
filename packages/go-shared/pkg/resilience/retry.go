package resilience

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)


type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      float64
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Jitter:      0.2,
	}
}

type RetryableFunc func(context.Context) error

func Retry(ctx context.Context, cfg RetryConfig, fn RetryableFunc) error {
	var err error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		delay := calcBackoff(cfg, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

func calcBackoff(cfg RetryConfig, attempt int) time.Duration {
	delay := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt))
	if delay > float64(cfg.MaxDelay) {
		delay = float64(cfg.MaxDelay)
	}
	if cfg.Jitter > 0 {
		jitter := (rand.Float64()*2 - 1) * cfg.Jitter * delay
		delay += jitter
		if delay < 0 {
			delay = 0
		}
	}
	return time.Duration(delay)
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	mu           sync.Mutex
	state        CircuitState
	failureCount int
	threshold    int
	resetTimeout time.Duration
	lastFailure  time.Time
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        CircuitClosed,
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	if err := cb.recordAttempt(); err != nil {
		return err
	}
	if err := fn(ctx); err != nil {
		cb.recordFailure()
		return err
	}
	cb.recordSuccess()
	return nil
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failureCount++
	cb.lastFailure = time.Now()
	if cb.failureCount >= cb.threshold {
		cb.state = CircuitOpen
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}
	cb.failureCount = 0
}

func (cb *CircuitBreaker) recordAttempt() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.maybeReset()
	if cb.state == CircuitOpen {
		return ErrCircuitOpen
	}
	return nil
}

func (cb *CircuitBreaker) maybeReset() {
	if cb.state == CircuitOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
		cb.state = CircuitHalfOpen
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.maybeReset()
	return cb.state
}
