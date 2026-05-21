package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/payment/internal/transport/http/middleware"
)

type Router struct {
	handler            *Handler
	authMw             gin.HandlerFunc
	webhookMiddlewares []gin.HandlerFunc
}

func NewRouter(handler *Handler, authMw gin.HandlerFunc, webhookMiddlewares ...gin.HandlerFunc) *Router {
	return &Router{handler: handler, authMw: authMw, webhookMiddlewares: webhookMiddlewares}
}

func (r *Router) Setup(engine *gin.Engine) {
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
	webhooks := v1.Group("/webhooks")
	for _, mw := range r.webhookMiddlewares {
		webhooks.Use(mw)
	}
	webhooks.POST("/:provider", r.handler.HandleWebhook)
}
