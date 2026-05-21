package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoWithRetry_Success(t *testing.T) {
	attempts := 0
	err := DoWithRetry(context.Background(), DefaultRetryConfig(), "test", func(ctx context.Context) error {
		attempts++
		return nil
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestDoWithRetry_RetriesOnFailure(t *testing.T) {
	attempts := 0
	cfg := DefaultRetryConfig()
	cfg.MaxAttempts = 3
	cfg.InitialInterval = 10 * time.Millisecond

	err := DoWithRetry(context.Background(), cfg, "test", func(ctx context.Context) error {
		attempts++
		return errors.New("connection refused")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDoWithRetry_NonRetryable(t *testing.T) {
	attempts := 0
	err := DoWithRetry(context.Background(), DefaultRetryConfig(), "test", func(ctx context.Context) error {
		attempts++
		return errors.New("invalid request")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestDoWithRetry_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := DoWithRetry(ctx, DefaultRetryConfig(), "test", func(ctx context.Context) error {
		return errors.New("connection refused")
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDefaultRetryableCheck(t *testing.T) {
	tests := []struct {
		err      error
		retryable bool
	}{
		{errors.New("connection refused"), true},
		{errors.New("connection reset by peer"), true},
		{errors.New("i/o timeout"), true},
		{errors.New("no such host"), true},
		{errors.New("circuit breaker is open"), true},
		{errors.New("invalid request"), false},
		{errors.New("bad gateway"), false},
	}

	for _, tt := range tests {
		result := DefaultRetryableCheck(tt.err)
		if result != tt.retryable {
			t.Errorf("DefaultRetryableCheck(%q) = %v, want %v", tt.err.Error(), result, tt.retryable)
		}
	}
}

func TestCalculateBackoff_Increases(t *testing.T) {
	cfg := DefaultRetryConfig()
	backoff1 := calculateBackoff(1, cfg)
	backoff2 := calculateBackoff(2, cfg)
	if backoff2 <= backoff1 {
		t.Errorf("expected backoff to increase, attempt 1: %v, attempt 2: %v", backoff1, backoff2)
	}
}

func TestCalculateBackoff_RespectsMax(t *testing.T) {
	cfg := DefaultRetryConfig()
	cfg.MaxInterval = 200 * time.Millisecond
	backoff := calculateBackoff(10, cfg)
	if backoff > cfg.MaxInterval {
		t.Errorf("backoff %v exceeds max %v", backoff, cfg.MaxInterval)
	}
}
