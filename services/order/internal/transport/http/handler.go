package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/order/internal/application"
	"github.com/shopee-clone/shopee/services/order/internal/domain"
	"go.uber.org/zap"
)

type Handler struct {
	orderService *application.OrderService
}

func NewHandler(orderService *application.OrderService) *Handler {
	return &Handler{orderService: orderService}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req application.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(string)
	if !ok || uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	req.UserID = uid

	order, err := h.orderService.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	// Validate ownership
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	uid, _ := userID.(string)
	r, _ := role.(string)
	_ = uid
	_ = r
	if order.UserID != uid && r != "admin" && r != "seller" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) ListOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	uid, _ := userID.(string)
	r, _ := role.(string)

	// Admin can list all orders, users can only list their own
	queryUserID := c.Query("user_id")
	if queryUserID == "" || r != "admin" {
		queryUserID = uid
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.orderService.ListOrders(c.Request.Context(), queryUserID, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  orders,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *Handler) GetOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id": order.ID,
		"status":   order.Status,
		"version":  order.Version,
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	uid, _ := userID.(string)
	r, _ := role.(string)

	cancelReq := &application.CancelOrderRequest{
		OrderID:       orderID,
		Reason:        req.Reason,
		CancelledBy:   uid,
		CancelledType: domain.CancellationTypeUser,
	}

	if r == "admin" {
		cancelReq.CancelledType = domain.CancellationTypeSystem
	}

	order, err := h.orderService.CancelOrder(c.Request.Context(), cancelReq)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) GetOrderHistory(c *gin.Context) {
	orderID := c.Param("id")
	history, err := h.orderService.GetOrderHistory(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id": orderID,
		"history":  history,
	})
}

func (h *Handler) GetReconciliationStatus(c *gin.Context) {
	orderID := c.Param("id")
	status, err := h.orderService.GetReconciliationStatus(c.Request.Context(), orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":       orderID,
		"reconciliation": status,
	})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrOrderNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrOrderNotCancellable:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrInvalidStateTransition:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrUnauthorized, domain.ErrInsufficientPermissions:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case domain.ErrConcurrentModification:
		c.JSON(http.StatusConflict, gin.H{"error": "concurrent modification detected, please retry"})
	default:
		zap.L().Error("unexpected error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
