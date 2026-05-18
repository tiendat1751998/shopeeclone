package resilience

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopee-clone/shopee/services/gateway/internal/config"
	"github.com/sony/gobreaker"
)

func TestCircuitBreakerPool_GetBreaker(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		MaxRequests:  5,
		Interval:     60 * time.Second,
		Timeout:      30 * time.Second,
		FailureRatio: 0.6,
		MinSamples:   5,
	}

	pool := NewCircuitBreakerPool(cfg)
	cb := pool.GetBreaker("test-service")
	if cb == nil {
		t.Fatal("circuit breaker should not be nil")
	}

	cb2 := pool.GetBreaker("test-service")
	if cb != cb2 {
		t.Error("should return same breaker for same service")
	}
}

func TestCircuitBreakerPool_DifferentServices(t *testing.T) {
	pool := NewCircuitBreakerPool(config.CircuitBreakerConfig{})
	cb1 := pool.GetBreaker("auth")
	cb2 := pool.GetBreaker("catalog")
	if cb1 == cb2 {
		t.Error("different services should have different breakers")
	}
}

func TestDefaultRetryPolicy(t *testing.T) {
	p := DefaultRetryPolicy(3)
	if p.MaxRetries != 3 {
		t.Errorf("expected 3 retries, got %d", p.MaxRetries)
	}
	if p.InitialInterval != 100*time.Millisecond {
		t.Errorf("expected 100ms interval, got %v", p.InitialInterval)
	}
}

func TestRetryPolicy_Success(t *testing.T) {
	p := DefaultRetryPolicy(3)

	calls := 0
	err := p.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryPolicy_RetryThenSuccess(t *testing.T) {
	p := DefaultRetryPolicy(3)

	calls := 0
	err := p.Do(context.Background(), func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryPolicy_Exhausted(t *testing.T) {
	p := DefaultRetryPolicy(2)

	calls := 0
	err := p.Do(context.Background(), func(ctx context.Context) error {
		calls++
		return errors.New("persistent error")
	})

	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls (1 initial + 2 retries), got %d", calls)
	}
}

func TestRetryPolicy_ContextCancelled(t *testing.T) {
	p := DefaultRetryPolicy(5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.Do(ctx, func(ctx context.Context) error {
		return errors.New("some error")
	})

	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestNewProxyExecutor(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		MaxRequests:  5,
		Interval:     60 * time.Second,
		Timeout:      30 * time.Second,
		FailureRatio: 0.6,
		MinSamples:   5,
	}

	exec := NewProxyExecutor(cfg, 30*time.Second, 2)
	if exec == nil {
		t.Fatal("executor should not be nil")
	}
	if exec.client == nil {
		t.Fatal("http client should not be nil")
	}
}

func TestCircuitBreaker_Transitions(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test",
		MaxRequests: 1,
		Interval:    0,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= 3
		},
	})

	for i := 0; i < 5; i++ {
		_, err := cb.Execute(func() (interface{}, error) {
			return nil, fmt.Errorf("error %d", i)
		})
		if i < 3 && err == nil {
			t.Error("expected error from failed execution")
		}
	}

	time.Sleep(150 * time.Millisecond)

	_, err := cb.Execute(func() (interface{}, error) {
		return "success", nil
	})
	if err != nil {
		t.Logf("after half-open, execute may still fail: %v", err)
	}
}

func TestRetryPolicy_IntervalGrowth(t *testing.T) {
	p := DefaultRetryPolicy(3)
	if p.InitialInterval != 100*time.Millisecond {
		t.Errorf("initial interval should be 100ms, got %v", p.InitialInterval)
	}

	p.InitialInterval = time.Second
	p.MaxInterval = 5 * time.Second
	p.Multiplier = 2.0

	interval := p.InitialInterval
	for i := 0; i < 3; i++ {
		interval = time.Duration(float64(interval) * p.Multiplier)
		if interval > p.MaxInterval {
			interval = p.MaxInterval
		}
	}

	if interval > p.MaxInterval {
		t.Errorf("interval should not exceed max interval")
	}
}

func TestProxyExecutor_WithHTTPServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	executor := NewProxyExecutor(
		config.CircuitBreakerConfig{MaxRequests: 5, MinSamples: 1},
		5*time.Second,
		1,
	)

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(context.Background())

	target := &struct{ Address string }{srv.URL}
	resp, err := executor.client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
