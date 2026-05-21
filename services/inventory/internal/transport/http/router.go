package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/inventory/internal/transport/http/middleware"
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

	inv := v1.Group("/inventory")
	if r.authMw != nil { inv.Use(r.authMw) }
	{
		inv.POST("/reserve", r.handler.ReserveStock)
		inv.POST("/release/:id", r.handler.ReleaseStock)
		inv.GET("/stock/:sku_id", r.handler.GetStock)
	}
}
