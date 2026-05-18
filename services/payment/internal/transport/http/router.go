package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/payment/internal/transport/http/middleware"
)

type Router struct {
	handler *Handler
	authMw  gin.HandlerFunc
}

func NewRouter(handler *Handler, authMw gin.HandlerFunc) *Router {
	return &Router{handler: handler, authMw: authMw}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.GET("/health/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "alive"}) })
	engine.GET("/health/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })

	v1 := engine.Group("/api/v1")
	v1.Use(middleware.RequestID())
	v1.Use(middleware.Recovery())

	payments := v1.Group("/payments")
	if r.authMw != nil { payments.Use(r.authMw) }
	{
		payments.POST("", r.handler.AuthorizePayment)
		payments.GET("/:id", r.handler.GetPayment)
		payments.POST("/:id/capture", r.handler.CapturePayment)
		payments.POST("/:id/refund", r.handler.RefundPayment)
	}

	// Webhook endpoint (no auth, uses signature verification)
	v1.POST("/webhooks/:provider", r.handler.HandleWebhook)
}
