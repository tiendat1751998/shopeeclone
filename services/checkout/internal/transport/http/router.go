package http

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/auth"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
)

type Router struct {
	handler   *Handler
	health    *health.Checker
	jwtSecret string
	redis     *redis.Client
}

func NewRouter(handler *Handler, healthChecker *health.Checker, jwtSecret string, redisClient *redis.Client) *Router {
	return &Router{handler: handler, health: healthChecker, jwtSecret: jwtSecret, redis: redisClient}
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

	// [SECURITY] Authentication is mandatory. If JWT secret is not configured,
	// we fail fast in production rather than allowing unauthenticated access.
	if r.jwtSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET is required - cannot start checkout service without authentication")
	}
	api.Use(auth.GinJWTAuth(r.jwtSecret))

	{
		api.POST("/checkout",
			middleware.RedisSlidingWindowLimiter(r.redis, 1, 5*time.Second, "checkout"),
			r.handler.InitiateCheckout,
		)
		api.GET("/checkout/:checkout_id/status", r.handler.GetCheckoutStatus)
		api.POST("/checkout/:checkout_id/retry", r.handler.RetryCheckout)
	}
}