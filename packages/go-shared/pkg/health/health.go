package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CheckFunc func(ctx context.Context) error

type Checker struct {
	mu       sync.RWMutex
	checks   map[string]CheckFunc
	started  time.Time
	version  string
	service  string
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Service   string            `json:"service"`
	Uptime    string            `json:"uptime"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

func NewChecker(service, version string) *Checker {
	return &Checker{
		checks:  make(map[string]CheckFunc),
		started: time.Now(),
		version: version,
		service: service,
	}
}

func (c *Checker) AddCheck(name string, fn CheckFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = fn
}

func (c *Checker) LivenessHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "alive",
			"service": c.service,
		})
	}
}

func (c *Checker) ReadinessHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		status := HealthStatus{
			Status:    "healthy",
			Version:   c.version,
			Service:   c.service,
			Uptime:    time.Since(c.started).String(),
			Timestamp: time.Now(),
			Checks:    make(map[string]string),
		}

		dCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
		defer cancel()

		healthy := true
		for name, check := range c.checks {
			if err := check(dCtx); err != nil {
				status.Checks[name] = err.Error()
				healthy = false
			} else {
				status.Checks[name] = "ok"
			}
		}

		if !healthy {
			status.Status = "unhealthy"
			ctx.JSON(http.StatusServiceUnavailable, status)
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}

type ProbeResponse struct {
	RequestID string `json:"request_id"`
	Healthy   bool   `json:"healthy"`
}

func LivenessProbe(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, ProbeResponse{
			RequestID: c.GetString("request_id"),
			Healthy:   true,
		})
	}
}

func ReadinessProbe(service string, checks ...CheckFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, check := range checks {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			if err := check(ctx); err != nil {
				cancel()
				c.JSON(http.StatusServiceUnavailable, ProbeResponse{
					RequestID: c.GetString("request_id"),
					Healthy:   false,
				})
				return
			}
			cancel()
		}

		c.JSON(http.StatusOK, ProbeResponse{
			RequestID: uuid.New().String(),
			Healthy:   true,
		})
	}
}
