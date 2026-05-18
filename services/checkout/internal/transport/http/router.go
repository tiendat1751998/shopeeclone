package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
)

type Router struct {
	handler *Handler
	health  *health.Checker
}

func NewRouter(handler *Handler, healthChecker *health.Checker) *Router {
	return &Router{handler: handler, health: healthChecker}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(
		middleware.Recovery(), middleware.ErrorHandler(),
		middleware.RequestID(), middleware.CORS(),
		middleware.OTelMiddleware("shopee-checkout"),
		observability.ObserveHTTPMetrics("shopee-checkout"),
	)

	engine.GET("/health", r.health.LivenessHandler())
	engine.GET("/ready", r.health.ReadinessHandler())
	engine.GET("/metrics", observability.MetricsHandler())

	api := engine.Group("/api/v1")
	{
		api.POST("/checkout", r.handler.InitiateCheckout)
		api.GET("/checkout/:checkout_id/status", r.handler.GetCheckoutStatus)
		api.POST("/checkout/:checkout_id/retry", r.handler.RetryCheckout)
	}
}
