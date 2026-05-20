package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/auth"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/product-catalog/internal/health"
)

type Router struct {
	handler   *Handler
	health    *health.Checker
	jwtSecret string
}

func NewRouter(h *Handler, hc *health.Checker, jwtSecret string) *Router {
	return &Router{handler: h, health: hc, jwtSecret: jwtSecret}
}

func (r *Router) Setup(e *gin.Engine) {
	e.Use(middleware.Recovery(), middleware.ErrorHandler(), middleware.RequestID(), middleware.CORS(), middleware.OTelMiddleware("shopee-product-catalog"), observability.ObserveHTTPMetrics("shopee-product-catalog"))
	e.GET("/health", r.health.LivenessHandler())
	e.GET("/ready", r.health.ReadinessHandler())
	e.GET("/metrics", observability.MetricsHandler())
	api := e.Group("/api/v1")
	if r.jwtSecret != "" {
		api.Use(auth.GinJWTAuth(r.jwtSecret))
	}
	{
		api.POST("/products", r.handler.CreateProduct)
		api.GET("/products/:id", r.handler.GetProduct)
		api.PUT("/products/:id", r.handler.UpdateProduct)
		api.DELETE("/products/:id", r.handler.ArchiveProduct)
		api.GET("/products/:product_id/skus", r.handler.AddSKU)
		api.POST("/products/:product_id/skus", r.handler.AddSKU)
		api.GET("/categories", r.handler.GetCategories)
		api.POST("/categories", r.handler.CreateCategory)
	}
}
