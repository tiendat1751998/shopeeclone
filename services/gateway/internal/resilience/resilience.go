package resilience

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sony/gobreaker"
	"github.com/shopee-clone/shopee/services/gateway/internal/config"
	"github.com/shopee-clone/shopee/services/gateway/internal/discovery"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CircuitBreakerPool struct {
	mu       sync.RWMutex
	breakers map[string]*gobreaker.CircuitBreaker
	cfg      config.CircuitBreakerConfig
}

func NewCircuitBreakerPool(cfg config.CircuitBreakerConfig) *CircuitBreakerPool {
	return &CircuitBreakerPool{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		cfg:      cfg,
	}
}

func (p *CircuitBreakerPool) GetBreaker(serviceName string) *gobreaker.CircuitBreaker {
	p.mu.RLock()
	cb, exists := p.breakers[serviceName]
	p.mu.RUnlock()

	if exists {
		return cb
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if cb, exists := p.breakers[serviceName]; exists {
		return cb
	}

	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        fmt.Sprintf("gateway-cb-%s", serviceName),
		MaxRequests: p.cfg.MaxRequests,
		Interval:    p.cfg.Interval,
		Timeout:     p.cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= p.cfg.MinSamples && failureRatio >= p.cfg.FailureRatio
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			observability.GetLogger().Info("circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	})

	p.breakers[serviceName] = cb
	return cb
}

type RetryPolicy struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	RetryableFunc   func(error) bool
}

func DefaultRetryPolicy(maxRetries int) *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:      maxRetries,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		RetryableFunc: func(err error) bool {
			if err == nil {
				return false
			}
			return true
		},
	}
}

func (p *RetryPolicy) Do(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error
	interval := p.InitialInterval

	for i := 0; i <= p.MaxRetries; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
			interval = time.Duration(float64(interval) * p.Multiplier)
			if interval > p.MaxInterval {
				interval = p.MaxInterval
			}
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err
		if !p.RetryableFunc(err) {
			return err
		}
		// Don't retry if context is done
		if ctx.Err() != nil {
			return ctx.Err()
		}

		observability.GetLogger().Warn("retrying request",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", p.MaxRetries),
			zap.Error(err),
		)
	}

	return fmt.Errorf("request failed after %d retries: %w", p.MaxRetries, lastErr)
}

type ProxyExecutor struct {
	cbPool       *CircuitBreakerPool
	retryPolicy  *RetryPolicy
	client       *http.Client
	timeout      time.Duration
}

func NewProxyExecutor(cfg config.CircuitBreakerConfig, timeout time.Duration, maxRetries int) *ProxyExecutor {
	return &ProxyExecutor{
		cbPool:      NewCircuitBreakerPool(cfg),
		retryPolicy: DefaultRetryPolicy(maxRetries),
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  false,
			},
		},
		timeout: timeout,
	}
}

type ProxyResponse struct {
	StatusCode int
	Header     http.Header
	Body       io.ReadCloser
}

func (e *ProxyExecutor) ExecuteWithBreaker(ctx context.Context, service string, target *discovery.ServiceTarget, req *http.Request) (*ProxyResponse, error) {
	ctx, span := otel.Tracer("shopee-gateway").Start(ctx, fmt.Sprintf("proxy.%s", service),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("upstream.service", service),
		attribute.String("upstream.address", target.Address),
	)

	cb := e.cbPool.GetBreaker(service)

	result, err := cb.Execute(func() (interface{}, error) {
		var resp *http.Response
		var err error

		err = e.retryPolicy.Do(ctx, func(rCtx context.Context) error {
			proxyReq, proxyErr := http.NewRequestWithContext(rCtx, req.Method, req.URL.String(), req.Body)
			if proxyErr != nil {
				return proxyErr
			}
			proxyReq.Header = req.Header.Clone()

			resp, err = e.client.Do(proxyReq)
			if err != nil {
				return err
			}

			if resp.StatusCode >= 500 {
				resp.Body.Close()
				return fmt.Errorf("upstream error: %d", resp.StatusCode)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		return &ProxyResponse{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       resp.Body,
		}, nil
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if err == gobreaker.ErrOpenState {
			observability.BusinessErrorsTotal.WithLabelValues("gateway", "CIRCUIT_BREAKER_OPEN").Inc()
			return nil, fmt.Errorf("circuit breaker open for service %s", service)
		}

		return nil, err
	}

	return result.(*ProxyResponse), nil
}
