package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		orders := api.Group("/orders")
		orders.POST("", h.CreateOrder)
		orders.GET("", h.ListOrders)
		orders.GET("/:id", h.GetOrder)
		orders.PUT("/:id/status", h.UpdateOrderStatus)

		inventoryGroup := api.Group("/inventory")
		inventoryGroup.POST("/reserve", h.ReserveInventory)
		inventoryGroup.POST("/release", h.ReleaseInventory)
		inventoryGroup.GET("/stock", h.CheckStock)

		picking := api.Group("/picking")
		picking.POST("/picklists", h.CreatePickList)
		picking.PUT("/picklists/:id/complete", h.CompletePicking)

		packing := api.Group("/packing")
		packing.POST("/create", h.CreatePacking)

		shipping := api.Group("/shipping")
		shipping.POST("/create", h.CreateShipment)

		returns := api.Group("/returns")
		returns.POST("", h.RequestReturn)
		returns.PUT("/:id/approve", h.ApproveReturn)
		returns.PUT("/:id/receive", h.ReceiveReturn)

		warehouses := api.Group("/warehouses")
		warehouses.GET("", h.ListWarehouses)
		warehouses.POST("/movement", h.RecordMovement)
	}
	return r
}
