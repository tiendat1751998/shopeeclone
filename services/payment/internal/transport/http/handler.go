package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/payment/internal/application"
	"github.com/shopee-clone/shopee/services/payment/internal/domain"
	"go.uber.org/zap"
)

type Handler struct {
	paymentService *application.PaymentService
}

func NewHandler(paymentService *application.PaymentService) *Handler {
	return &Handler{paymentService: paymentService}
}

func (h *Handler) AuthorizePayment(c *gin.Context) {
	var req application.AuthorizePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	req.UserID = userID.(string)
	payment, err := h.paymentService.AuthorizePayment(c.Request.Context(), &req)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, payment)
}

func (h *Handler) CapturePayment(c *gin.Context) {
	paymentID := c.Param("id")
	userID, _ := c.Get("user_id")
	payment, err := h.paymentService.CapturePayment(c.Request.Context(), paymentID, userID.(string))
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, payment)
}

func (h *Handler) RefundPayment(c *gin.Context) {
	paymentID := c.Param("id")
	var req struct {
		Reason         string `json:"reason" binding:"required"`
		Amount         int64  `json:"amount" binding:"required"`
		IdempotencyKey string `json:"idempotency_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refund, err := h.paymentService.RefundPayment(c.Request.Context(), paymentID, req.Reason, req.IdempotencyKey, req.Amount)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, refund)
}

func (h *Handler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	payment, err := h.paymentService.GetPayment(c.Request.Context(), paymentID)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, payment)
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	pspProvider := c.Param("provider")
	eventType := c.Query("event_type")
	signature := c.GetHeader("X-Webhook-Signature")
	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	payload, _ := c.GetRawData()
	if err := h.paymentService.HandleWebhook(c.Request.Context(), pspProvider, eventType, payload, signature, idempotencyKey); err != nil {
		handleError(c, err); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrPaymentNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrDoubleChargeDetected:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrInvalidWebhookSignature:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case domain.ErrWebhookReplayDetected:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrRefundNotAllowed:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		zap.L().Error("unexpected error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
