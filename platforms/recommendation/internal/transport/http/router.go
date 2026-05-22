package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/health"
)

type Router struct {
	handler *Handler
	health  *health.Checker
}

func NewRouter(h *Handler, hc *health.Checker) *Router {
	return &Router{handler: h, health: hc}
}

func (r *Router) Setup(e *gin.Engine) {
	e.Use(middleware.Recovery(), middleware.ErrorHandler(), middleware.RequestID(), middleware.CORS(), middleware.OTelMiddleware("shopee-recommendation"), observability.ObserveHTTPMetrics("shopee-recommendation"))

	e.GET("/health", r.health.LivenessHandler())
	e.GET("/ready", r.health.ReadinessHandler())
	e.GET("/metrics", observability.MetricsHandler())

	{
		e.POST("/", r.handler.GetRecommendations)
		e.POST("/batch", r.handler.BatchRecommendations)
		e.POST("/feedback", r.handler.RecordFeedback)
		e.GET("/trending", r.handler.GetTrending)
	}
}
