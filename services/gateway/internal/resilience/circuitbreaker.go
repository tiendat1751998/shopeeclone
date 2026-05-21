package resilience

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type State int32

const (
	StateClosed   State = 0
	StateHalfOpen State = 1
	StateOpen     State = 2
)

type CircuitBreaker struct {
	name          string
	maxRequests   int32
	interval      time.Duration
	timeout       time.Duration
	failureRatio  float64
	minSamples    int32

	state     int32
	counts    Counts
	expiry    time.Time
	lastState time.Time

	mu sync.RWMutex
}

type Counts struct {
	Requests             int32
	TotalSuccesses       int32
	TotalFailures        int32
	ConsecutiveSuccesses int32
	ConsecutiveFailures  int32
}

func NewCircuitBreaker(name string, opts CircuitBreakerOptions) *CircuitBreaker {
	if opts.MaxRequests <= 0 {
		opts.MaxRequests = 5
	}
	if opts.Interval <= 0 {
		opts.Interval = 60 * time.Second
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.FailureRatio <= 0 {
		opts.FailureRatio = 0.6
	}
	if opts.MinSamples <= 0 {
		opts.MinSamples = 5
	}
	return &CircuitBreaker{
		name:          name,
		maxRequests:   int32(opts.MaxRequests),
		interval:      opts.Interval,
		timeout:       opts.Timeout,
		failureRatio:  opts.FailureRatio,
		minSamples:    int32(opts.MinSamples),
		state:         int32(StateClosed),
		lastState:     time.Now(),
	}
}

type CircuitBreakerOptions struct {
	MaxRequests  int
	Interval     time.Duration
	Timeout      time.Duration
	FailureRatio float64
	MinSamples   int
}

func (cb *CircuitBreaker) Name() string { return cb.name }

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allow() {
		return fmt.Errorf("circuit breaker %s is open", cb.name)
	}

	err := fn()

	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}

	return err
}

func (cb *CircuitBreaker) allow() bool {
	state := State(atomic.LoadInt32(&cb.state))

	if state == StateHalfOpen {
		return true
	}

	if state == StateOpen {
		cb.mu.RLock()
		expiry := cb.expiry
		cb.mu.RUnlock()

		if time.Now().After(expiry) {
			cb.setState(StateHalfOpen)
			return true
		}
		return false
	}

	cb.mu.RLock()
	counts := cb.counts
	cb.mu.RUnlock()

	if counts.Requests < cb.minSamples {
		return true
	}

	failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
	if failureRate >= cb.failureRatio {
		cb.setState(StateOpen)
		return false
	}

	return true
}

func (cb *CircuitBreaker) recordSuccess() {
	state := State(atomic.LoadInt32(&cb.state))

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.counts.Requests++
	cb.counts.TotalSuccesses++
	cb.counts.ConsecutiveSuccesses++
	cb.counts.ConsecutiveFailures = 0

	if state == StateHalfOpen && cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
		cb.resetCounts()
		atomic.StoreInt32(&cb.state, int32(StateClosed))
		cb.lastState = time.Now()
	}
}

func (cb *CircuitBreaker) recordFailure() {
	state := State(atomic.LoadInt32(&cb.state))

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.counts.Requests++
	cb.counts.TotalFailures++
	cb.counts.ConsecutiveFailures++
	cb.counts.ConsecutiveSuccesses = 0

	if state == StateHalfOpen {
		atomic.StoreInt32(&cb.state, int32(StateOpen))
		cb.expiry = time.Now().Add(cb.timeout)
		cb.lastState = time.Now()
		return
	}

	if cb.counts.Requests >= cb.minSamples {
		failureRate := float64(cb.counts.TotalFailures) / float64(cb.counts.Requests)
		if failureRate >= cb.failureRatio {
			atomic.StoreInt32(&cb.state, int32(StateOpen))
			cb.expiry = time.Now().Add(cb.timeout)
			cb.lastState = time.Now()
		}
	}
}

func (cb *CircuitBreaker) setState(state State) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	atomic.StoreInt32(&cb.state, int32(state))
	cb.lastState = time.Now()

	if state == StateOpen {
		cb.expiry = time.Now().Add(cb.timeout)
	} else if state == StateClosed {
		cb.resetCounts()
	}
}

func (cb *CircuitBreaker) resetCounts() {
	cb.counts = Counts{}
}

func (cb *CircuitBreaker) State() State {
	return State(atomic.LoadInt32(&cb.state))
}

func (cb *CircuitBreaker) Counts() Counts {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.counts
}
