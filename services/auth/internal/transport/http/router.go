package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/health"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
)

type Router struct {
	handler *Handler
	health  *health.Checker
	redis   *redis.Client
}

func NewRouter(handler *Handler, healthChecker *health.Checker, redisClient *redis.Client) *Router {
	return &Router{
		handler: handler,
		health:  healthChecker,
		redis:   redisClient,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(
		middleware.Recovery(),
		middleware.ErrorHandler(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.OTelMiddleware("shopee-auth"),
		observability.ObserveHTTPMetrics("shopee-auth"),
	)

	engine.GET("/health", r.health.LivenessHandler())
	engine.GET("/ready", r.health.ReadinessHandler())
	engine.GET("/metrics", observability.MetricsHandler())

	api := engine.Group("/api/v1/auth")
	{
		api.POST("/register", r.handler.Register)
		api.POST("/login",
			middleware.RedisSlidingWindowLimiter(r.redis, 5, 5*time.Minute, "auth_login"),
			r.handler.Login,
		)
		api.POST("/refresh", r.handler.RefreshToken)
		api.POST("/logout", r.handler.Logout)
		api.POST("/logout/all", r.handler.LogoutAll)
		api.GET("/sessions", r.handler.GetSessions)
		api.DELETE("/sessions/:session_id", r.handler.RevokeSession)
		api.GET("/profile", r.handler.GetProfile)
		api.POST("/validate", r.handler.ValidateToken)
	}
}
