package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/gateway/internal/discovery"
	"github.com/shopee-clone/shopee/services/gateway/internal/middleware"
	"github.com/shopee-clone/shopee/services/gateway/internal/resilience"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Proxy struct {
	executor   *resilience.ProxyExecutor
	discovery  *discovery.ServiceDiscovery
	transport  *http.Transport
}

func NewProxy(
	executor *resilience.ProxyExecutor,
	svcDiscovery *discovery.ServiceDiscovery,
	maxIdleConns int,
	idleConnTimeout time.Duration,
) *Proxy {
	return &Proxy{
		executor:  executor,
		discovery: svcDiscovery,
		transport: &http.Transport{
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConns / 4,
			IdleConnTimeout:     idleConnTimeout,
			DisableCompression:  false,
		},
	}
}

type ProxyTarget struct {
	ServiceName string
	PathPrefix  string
	StripPrefix string
	Timeout     time.Duration
}

func (p *Proxy) ReverseProxy(target *ProxyTarget) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, span := otel.Tracer("shopee-gateway").Start(c.Request.Context(),
			fmt.Sprintf("reverse_proxy.%s", target.ServiceName),
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		span.SetAttributes(
			attribute.String("upstream.service", target.ServiceName),
			attribute.String("upstream.path", c.Request.URL.Path),
		)

		instance, err := p.discovery.GetInstance(target.ServiceName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "no healthy upstream instances")

			observability.GetLogger().Error("no healthy upstream instances",
				zap.String("service", target.ServiceName),
				zap.Error(err),
			)

			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error_code": "SERVICE_UNAVAILABLE",
				"message":    fmt.Sprintf("service %s is currently unavailable", target.ServiceName),
			})
			return
		}

		targetURL := fmt.Sprintf("http://%s:%d", instance.Address, instance.Port)
		upstreamURL, err := url.Parse(targetURL)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid upstream url"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
		proxy.Transport = p.transport
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			observability.GetLogger().Error("proxy error",
				zap.String("service", target.ServiceName),
				zap.Error(err),
			)
			observability.BusinessErrorsTotal.WithLabelValues("gateway", "PROXY_ERROR").Inc()
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(gin.H{
				"error_code": "BAD_GATEWAY",
				"message":    "upstream connection error",
			})
		}

		p.modifyRequest(c, proxy, target)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *Proxy) modifyRequest(c *gin.Context, proxy *httputil.ReverseProxy, target *ProxyTarget) {
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)

		if target.StripPrefix != "" {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, target.StripPrefix)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}

		req.Host = req.URL.Host

		p.injectGatewayHeaders(c, req)
	}
}

func (p *Proxy) injectGatewayHeaders(c *gin.Context, req *http.Request) {
	if correlationID, exists := c.Get(string(middleware.CorrelationIDKey)); exists {
		req.Header.Set("X-Correlation-ID", fmt.Sprintf("%v", correlationID))
	}
	if requestID, exists := c.Get(string(middleware.RequestIDKey)); exists {
		req.Header.Set("X-Request-ID", fmt.Sprintf("%v", requestID))
	}
	if userID, exists := c.Get(string(middleware.UserIDKey)); exists {
		req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
	}
	if roles, exists := c.Get(string(middleware.UserRolesKey)); exists {
		if roleList, ok := roles.([]string); ok {
			req.Header.Set("X-User-Roles", strings.Join(roleList, ","))
		}
	}
	if deviceInfo, exists := c.Get(string(middleware.DeviceInfoKey)); exists {
		if info, ok := deviceInfo.(map[string]string); ok {
			for k, v := range info {
				req.Header.Set(fmt.Sprintf("X-Device-%s", k), v)
			}
		}
	}

	req.Header.Set("X-Forwarded-For", c.ClientIP())
	req.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)
	req.Header.Set("X-Forwarded-Host", c.Request.Host)

	traceID := trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()
	if traceID != "" {
		req.Header.Set("X-Trace-ID", traceID)
	}
}

func (p *Proxy) HealthCheck(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance, err := p.discovery.GetInstance(serviceName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"service": serviceName,
				"healthy": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"healthy": true,
			"address": fmt.Sprintf("%s:%d", instance.Address, instance.Port),
		})
	}
}

func copyRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func httpClientTimeout(serviceName string) time.Duration {
	switch serviceName {
	case "auth":
		return 10 * time.Second
	case "catalog":
		return 15 * time.Second
	case "cart":
		return 5 * time.Second
	case "order":
		return 30 * time.Second
	case "inventory":
		return 5 * time.Second
	case "payment":
		return 30 * time.Second
	case "search":
		return 10 * time.Second
	default:
		return 30 * time.Second
	}
}
