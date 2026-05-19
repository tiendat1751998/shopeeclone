package resilience

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type RetryPolicy struct {
	MaxRetries       int     `json:"max_retries"`
	BackoffInitialMs int     `json:"backoff_initial_ms"`
	BackoffMultiplier float64 `json:"backoff_multiplier"`
	MaxBackoffMs     int     `json:"max_backoff_ms"`
	RetryOn          []string `json:"retry_on"`
}

type TimeoutPolicy struct {
	RequestTimeoutMs int `json:"request_timeout_ms"`
	IdleTimeoutMs    int `json:"idle_timeout_ms"`
}

type CircuitState int32

const (
	StateClosed   CircuitState = 0
	StateOpen     CircuitState = 1
	StateHalfOpen CircuitState = 2
)

type CircuitBreaker struct {
	name          string
	state         atomic.Int32
	failureCount  atomic.Int64
	successCount  atomic.Int64
	threshold     int64
	halfOpenMax   int64
	recoveryTime  time.Duration
	lastStateChange time.Time
	mu            sync.Mutex
}

type Bulkhead struct {
	name          string
	maxConcurrent int64
	queueSize     int64
	sem           chan struct{}
}

func (cb *CircuitBreaker) Name() string {
	return cb.name
}

func NewCircuitBreaker(name string, threshold int64, halfOpenMax int64, recoveryTime time.Duration) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:         name,
		threshold:    threshold,
		halfOpenMax:  halfOpenMax,
		recoveryTime: recoveryTime,
	}
	cb.state.Store(int32(StateClosed))
	cb.lastStateChange = time.Now()
	return cb
}

func (cb *CircuitBreaker) State() CircuitState {
	return CircuitState(cb.state.Load())
}

func (cb *CircuitBreaker) AllowRequest() bool {
	state := cb.State()
	switch state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateChange) >= cb.recoveryTime {
			cb.mu.Lock()
			if cb.State() == StateOpen {
				cb.state.Store(int32(StateHalfOpen))
				cb.lastStateChange = time.Now()
			}
			cb.mu.Unlock()
			return true
		}
		return false
	case StateHalfOpen:
		success := cb.successCount.Load()
		return success < cb.halfOpenMax
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.successCount.Add(1)
	if cb.State() == StateHalfOpen && cb.successCount.Load() >= cb.halfOpenMax {
		cb.mu.Lock()
		cb.state.Store(int32(StateClosed))
		cb.failureCount.Store(0)
		cb.successCount.Store(0)
		cb.lastStateChange = time.Now()
		cb.mu.Unlock()
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.failureCount.Add(1)
	if cb.State() == StateHalfOpen {
		cb.mu.Lock()
		cb.state.Store(int32(StateOpen))
		cb.lastStateChange = time.Now()
		cb.mu.Unlock()
		return
	}
	if cb.State() == StateClosed && cb.failureCount.Load() >= cb.threshold {
		cb.mu.Lock()
		cb.state.Store(int32(StateOpen))
		cb.lastStateChange = time.Now()
		cb.mu.Unlock()
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state.Store(int32(StateClosed))
	cb.failureCount.Store(0)
	cb.successCount.Store(0)
	cb.lastStateChange = time.Now()
}

func NewBulkhead(name string, maxConcurrent, queueSize int64) *Bulkhead {
	return &Bulkhead{
		name:          name,
		maxConcurrent: maxConcurrent,
		queueSize:     queueSize,
		sem:           make(chan struct{}, maxConcurrent+queueSize),
	}
}

func (b *Bulkhead) Acquire(ctx context.Context) error {
	select {
	case b.sem <- struct{}{}:
		return nil
	default:
	}

	if int64(len(b.sem)) >= b.maxConcurrent+b.queueSize {
		return errors.New("bulkhead queue full")
	}

	select {
	case b.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Bulkhead) Release() {
	<-b.sem
}

type Executor struct {
	circuitBreakers map[string]*CircuitBreaker
	bulkheads       map[string]*Bulkhead
	mu              sync.RWMutex
}

func NewExecutor() *Executor {
	return &Executor{
		circuitBreakers: make(map[string]*CircuitBreaker),
		bulkheads:       make(map[string]*Bulkhead),
	}
}

func (e *Executor) AddCircuitBreaker(name string, threshold int64, halfOpenMax int64, recoveryTime time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.circuitBreakers[name] = NewCircuitBreaker(name, threshold, halfOpenMax, recoveryTime)
}

func (e *Executor) AddBulkhead(name string, maxConcurrent, queueSize int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.bulkheads[name] = NewBulkhead(name, maxConcurrent, queueSize)
}

func (e *Executor) GetCircuitBreaker(name string) *CircuitBreaker {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.circuitBreakers[name]
}

func (e *Executor) GetBulkhead(name string) *Bulkhead {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.bulkheads[name]
}

func (e *Executor) ListCircuitBreakers() []*CircuitBreaker {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var result []*CircuitBreaker
	for _, cb := range e.circuitBreakers {
		result = append(result, cb)
	}
	return result
}

func (e *Executor) ExecuteWithRetry(ctx context.Context, policy RetryPolicy, fn func(context.Context) error) error {
	var lastErr error
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := float64(policy.BackoffInitialMs) * math.Pow(policy.BackoffMultiplier, float64(attempt-1))
			if backoff > float64(policy.MaxBackoffMs) {
				backoff = float64(policy.MaxBackoffMs)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(backoff) * time.Millisecond):
			}
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		shouldRetry := false
		for _, retryOn := range policy.RetryOn {
			if err.Error() == retryOn || containsSubstring(err.Error(), retryOn) {
				shouldRetry = true
				break
			}
		}
		if !shouldRetry {
			return err
		}
	}
	return lastErr
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (e *Executor) ExecuteWithTimeout(ctx context.Context, policy TimeoutPolicy, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(policy.RequestTimeoutMs)*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(ctx)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *Executor) ExecuteWithBulkhead(ctx context.Context, name string, fn func(context.Context) error) error {
	bh := e.GetBulkhead(name)
	if bh == nil {
		return errors.New("bulkhead not found: " + name)
	}

	if err := bh.Acquire(ctx); err != nil {
		return err
	}
	defer bh.Release()

	return fn(ctx)
}
