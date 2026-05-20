package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/cart/internal/config"
	cartMiddleware "github.com/shopee-clone/shopee/services/cart/internal/transport/http/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
)

type Router struct {
	handler *Handler
	health  *health.Checker
	cfg     config.JWTConfig
}

func NewRouter(handler *Handler, healthChecker *health.Checker, jwtCfg config.JWTConfig) *Router {
	return &Router{handler: handler, health: healthChecker, cfg: jwtCfg}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(
		middleware.Recovery(),
		middleware.ErrorHandler(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.OTelMiddleware("shopee-cart"),
		observability.ObserveHTTPMetrics("shopee-cart"),
	)

	engine.GET("/health", r.health.LivenessHandler())
	engine.GET("/ready", r.health.ReadinessHandler())
	engine.GET("/metrics", observability.MetricsHandler())

	api := engine.Group("/api/v1")
	api.Use(cartMiddleware.JWTAuth(r.cfg))
	{
		api.GET("/carts/:cart_id", r.handler.GetCart)
		api.POST("/carts/:cart_id/items", r.handler.AddItem)
		api.PUT("/carts/:cart_id/items/:item_id", r.handler.UpdateItem)
		api.DELETE("/carts/:cart_id/items/:item_id", r.handler.RemoveItem)
		api.DELETE("/carts/:cart_id/items", r.handler.ClearCart)
		api.POST("/carts/merge", r.handler.MergeCarts)
		api.POST("/carts/:cart_id/checkout-preview", r.handler.CheckoutPreview)
	}
}
