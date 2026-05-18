package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/inventory/internal/application"
	"github.com/shopee-clone/shopee/services/inventory/internal/domain"
	"go.uber.org/zap"
)

type Handler struct {
	inventoryService *application.InventoryService
}

func NewHandler(svc *application.InventoryService) *Handler { return &Handler{inventoryService: svc} }

func (h *Handler) ReserveStock(c *gin.Context) {
	var req application.ReserveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	userID, _ := c.Get("user_id"); req.UserID = userID.(string)
	reservation, err := h.inventoryService.ReserveStock(c.Request.Context(), &req)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, reservation)
}

func (h *Handler) ReleaseStock(c *gin.Context) {
	reservationID := c.Param("id")
	if err := h.inventoryService.ReleaseStock(c.Request.Context(), reservationID); err != nil {
		handleError(c, err); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "released"})
}

func (h *Handler) GetStock(c *gin.Context) {
	skuID := c.Param("sku_id")
	warehouseID := c.Query("warehouse_id")
	stock, err := h.inventoryService.GetStock(c.Request.Context(), skuID, warehouseID)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, stock)
}

func handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrStockNotFound: c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrInsufficientStock: c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrReservationNotFound: c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrOversellPrevented: c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default: zap.L().Error("unexpected error", zap.Error(err)); c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
