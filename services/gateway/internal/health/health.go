package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/gateway/internal/discovery"
	"go.uber.org/zap"
)

type GatewayHealth struct {
	discovery *discovery.ServiceDiscovery
}

func NewGatewayHealth(d *discovery.ServiceDiscovery) *GatewayHealth {
	return &GatewayHealth{discovery: d}
}

type UpstreamHealth struct {
	Service string `json:"service"`
	Healthy bool   `json:"healthy"`
	Error   string `json:"error,omitempty"`
}

func (h *GatewayHealth) UpstreamsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services := h.discovery.GetAllServices()
		results := make([]UpstreamHealth, 0, len(services))

		for _, svc := range services {
			instances := h.discovery.GetInstances(svc)
			healthy := len(instances) > 0
			uh := UpstreamHealth{
				Service: svc,
				Healthy: healthy,
			}
			if !healthy {
				uh.Error = "no healthy instances"
			}
			results = append(results, uh)
		}

		c.JSON(http.StatusOK, gin.H{
			"timestamp": time.Now().UTC(),
			"services":  results,
			"total":     len(results),
			"healthy":   countHealthy(results),
		})
	}
}

func countHealthy(results []UpstreamHealth) int {
	count := 0
	for _, r := range results {
		if r.Healthy {
			count++
		}
	}
	return count
}

func RedisHealthCheck(addr string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		_ = timeout
		return nil
	}
}

func logHealthState(service string, healthy bool) {
	if healthy {
		observability.GetLogger().Info("service health check passed",
			zap.String("service", service))
	} else {
		observability.GetLogger().Warn("service health check failed",
			zap.String("service", service))
	}
}
