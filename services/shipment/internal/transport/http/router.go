package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/shipment/internal/transport/http/middleware"
)

type Router struct {
	handler *Handler
	authMw  gin.HandlerFunc
}

func NewRouter(handler *Handler, authMw gin.HandlerFunc) *Router { return &Router{handler: handler, authMw: authMw} }

func (r *Router) Setup(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	v1.Use(middleware.RequestID())
	v1.Use(middleware.Recovery())

	shipments := v1.Group("/shipments")
	if r.authMw != nil { shipments.Use(r.authMw) }
	{
		shipments.POST("", r.handler.CreateShipment)
		shipments.GET("/:id", r.handler.GetShipment)
		shipments.POST("/:id/status", r.handler.UpdateStatus)
		shipments.GET("/:id/tracking", r.handler.GetTracking)
	}

	v1.POST("/webhooks/:provider", r.handler.HandleWebhook)
}
