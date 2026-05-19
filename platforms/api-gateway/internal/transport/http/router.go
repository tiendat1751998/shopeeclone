package http

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(h *Handler) *Router {
	return &Router{handler: h}
}

func (r *Router) Setup(e *gin.Engine) {
	api := e.Group("/api/v1")
	{
		routes := api.Group("/routes")
		{
			routes.POST("", r.handler.RegisterRoute)
			routes.GET("", r.handler.ListRoutes)
			routes.DELETE("/:id", r.handler.DeregisterRoute)
			routes.POST("/match", r.handler.MatchRoute)
		}

		rateLimit := api.Group("/rate-limit")
		{
			rateLimit.POST("/rules", r.handler.CreateRateLimitRule)
			rateLimit.POST("/check", r.handler.CheckRateLimit)
		}

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/api-keys", r.handler.CreateAPIKey)
			authGroup.POST("/api-keys/validate", r.handler.ValidateAPIKey)
			authGroup.POST("/jwt/sign", r.handler.SignJWT)
			authGroup.POST("/jwt/verify", r.handler.VerifyJWT)
		}

		transforms := api.Group("/transforms")
		{
			transforms.POST("", r.handler.CreateTransformRule)
			transforms.POST("/apply", r.handler.ApplyTransform)
		}

		cb := api.Group("/circuit-breakers")
		{
			cb.POST("", r.handler.CreateCircuitBreaker)
			cb.GET("", r.handler.ListCircuitBreakers)
			cb.POST("/:id/record-success", r.handler.RecordCircuitBreakerSuccess)
			cb.POST("/:id/record-failure", r.handler.RecordCircuitBreakerFailure)
		}

		cache := api.Group("/edge-cache")
		{
			cache.GET("/:key", r.handler.GetCacheValue)
			cache.POST("", r.handler.SetCacheValue)
			cache.POST("/purge", r.handler.PurgeCache)
			cache.GET("/stats", r.handler.GetCacheStats)
		}
	}
}
