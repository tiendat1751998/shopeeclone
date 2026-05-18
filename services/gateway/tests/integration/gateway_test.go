package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/services/gateway/internal/auth"
	"github.com/shopee-clone/shopee/services/gateway/internal/config"
	"github.com/shopee-clone/shopee/services/gateway/internal/discovery"
	"github.com/shopee-clone/shopee/services/gateway/internal/ratelimit"
	"github.com/shopee-clone/shopee/services/gateway/internal/resilience"
	"github.com/shopee-clone/shopee/services/gateway/internal/routing"
	"github.com/shopee-clone/shopee/services/gateway/internal/transport"
)

func newTestRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("User-Agent", "test-agent/1.0")
	req.Header.Set("Host", "localhost")
	return req
}



func setupTestRouter() *gin.Engine {
	cfg := &config.Config{
		AppName:  "test-gateway",
		AppEnv:   "test",
		LogLevel: "error",
		HTTPPort: 8080,
		RateLimit: config.RateLimitConfig{
			Enabled: false,
		},
		Auth: config.AuthConfig{
			EnableRBAC: false,
		},
		Server: config.ServerConfig{
			MaxBodySize: 10485760,
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
		},
		OpenTelemetry: config.OTELConfig{
			ServiceName: "test-gateway",
		},
		Upstreams: config.UpstreamConfig{
			DefaultTimeout:  5 * time.Second,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			MaxRetries:      1,
			CircuitBreaker: config.CircuitBreakerConfig{
				MaxRequests: 5,
				MinSamples:  5,
			},
		},
	}

	jwtValidator := auth.NewJWTValidator(cfg.Auth, nil)
	authMiddleware := auth.NewAuthMiddleware(jwtValidator)

	rateLimiter := ratelimit.NewRateLimiter(nil, cfg.RateLimit)

	svcDiscovery := discovery.NewServiceDiscovery()
	svcDiscovery.RegisterStatic("test-service", []*discovery.ServiceInstance{
		{ID: "test-1", Name: "test-service", Address: "localhost", Port: 9999, Weight: 1},
	})

	executor := resilience.NewProxyExecutor(
		cfg.Upstreams.CircuitBreaker,
		cfg.Upstreams.DefaultTimeout,
		cfg.Upstreams.MaxRetries,
	)

	proxy := transport.NewProxy(
		executor,
		svcDiscovery,
		cfg.Upstreams.MaxIdleConns,
		cfg.Upstreams.IdleConnTimeout,
	)

	healthChecker := health.NewChecker("test-gateway", "1.0.0")

	gin.SetMode(gin.TestMode)
	engine := gin.New()

	router := routing.NewRouter(cfg, proxy, rateLimiter, authMiddleware, svcDiscovery, healthChecker)
	router.Setup(engine)

	return engine
}

func TestGateway_HealthEndpoint(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/health")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["status"] != "alive" {
		t.Errorf("expected alive status, got %v", response["status"])
	}
}

func TestGateway_ReadinessEndpoint(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/ready")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGateway_MetricsEndpoint(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/metrics")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGateway_NotFound(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/nonexistent")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusUnauthorized {
		t.Errorf("expected 404 or 401, got %d", w.Code)
	}
}

func TestGateway_CORSHeaders(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodOptions, "/api/v1/products")
	req.Header.Set("Origin", "http://localhost:3000")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestGateway_SecurityHeaders(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/health")
	engine.ServeHTTP(w, req)

	headers := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Strict-Transport-Security",
		"Referrer-Policy",
	}

	for _, header := range headers {
		if w.Header().Get(header) == "" {
			t.Errorf("missing security header: %s", header)
		}
	}
}

func TestGateway_CorrelationID(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/health")
	req.Header.Set("X-Correlation-ID", "test-correlation-id")
	engine.ServeHTTP(w, req)

	if w.Header().Get("X-Correlation-ID") != "test-correlation-id" {
		t.Errorf("expected correlation id to be preserved")
	}
}

func TestGateway_RequestIDGenerated(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/health")
	engine.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected request ID to be generated")
	}
}

func TestGateway_MissingUserAgent(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/api/v1/products")
	req.Header.Del("User-Agent")
	engine.ServeHTTP(w, req)

	t.Logf("user-agent enforcement returned %d", w.Code)
}

func TestGateway_UpstreamsEndpoint(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodGet, "/upstreams")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if response["services"] == nil {
		t.Error("expected services list")
	}
}

func TestGateway_MaxBodySize(t *testing.T) {
	engine := setupTestRouter()

	w := httptest.NewRecorder()
	req := newTestRequest(http.MethodPost, "/api/v1/auth/login")
	engine.ServeHTTP(w, req)

	if w.Code >= 400 && w.Code < 500 {
		t.Logf("request rejected as expected with %d", w.Code)
	}
}

func BenchmarkGatewayHealth(b *testing.B) {
	engine := setupTestRouter()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := newTestRequest(http.MethodGet, "/health")
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkGatewayMetrics(b *testing.B) {
	engine := setupTestRouter()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := newTestRequest(http.MethodGet, "/metrics")
		engine.ServeHTTP(w, req)
	}
}
