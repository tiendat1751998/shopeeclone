package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/order/internal/transport/http/middleware"
)

type Router struct {
	handler *Handler
	authMw  gin.HandlerFunc
}

func NewRouter(handler *Handler, authMw gin.HandlerFunc) *Router {
	return &Router{
		handler: handler,
		authMw:  authMw,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	// API v1
	v1 := engine.Group("/api/v1")
	v1.Use(middleware.RequestID())
	v1.Use(middleware.Logger())
	v1.Use(middleware.Recovery())

	// Protected routes
	orders := v1.Group("/orders")
	if r.authMw != nil {
		orders.Use(r.authMw)
	}
	{
		orders.POST("", r.handler.CreateOrder)
		orders.GET("", r.handler.ListOrders)
		orders.GET("/:id", r.handler.GetOrder)
		orders.GET("/:id/status", r.handler.GetOrderStatus)
		orders.POST("/:id/cancel", r.handler.CancelOrder)
		orders.GET("/:id/history", r.handler.GetOrderHistory)
		orders.GET("/:id/reconciliation", r.handler.GetReconciliationStatus)
	}
}
