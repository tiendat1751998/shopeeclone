package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/services/cart/internal/application"
	"github.com/shopee-clone/shopee/services/cart/internal/domain"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Handler struct {
	service *application.CartService
}

func NewHandler(service *application.CartService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCart(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.get_cart")
	defer span.End()

	cartID := c.Param("cart_id")
	cart, items, err := h.service.GetCartWithItems(ctx, cartID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart, "items": items})
}

func (h *Handler) AddItem(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.add_item")
	defer span.End()

	cartID := c.Param("cart_id")

	var req struct {
		SKU         string `json:"sku" binding:"required"`
		ProductName string `json:"product_name" binding:"required"`
		ShopID      string `json:"shop_id" binding:"required"`
		ShopName    string `json:"shop_name" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required,min=1"`
		UnitPrice   int64  `json:"unit_price" binding:"required,min=0"`
		ImageURL    string `json:"image_url"`
		Attributes  string `json:"attributes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	item, err := h.service.AddItem(ctx, cartID, application.AddItemRequest{
		SKU: req.SKU, ProductName: req.ProductName, ShopID: req.ShopID,
		ShopName: req.ShopName, Quantity: req.Quantity, UnitPrice: req.UnitPrice,
		ImageURL: req.ImageURL, Attributes: req.Attributes,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	span.SetAttributes(attribute.String("item_id", item.ID))
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.update_item")
	defer span.End()

	cartID := c.Param("cart_id")
	itemID := c.Param("item_id")

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	if err := h.service.UpdateItemQuantity(ctx, cartID, itemID, req.Quantity); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item updated"})
}

func (h *Handler) RemoveItem(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.remove_item")
	defer span.End()

	cartID := c.Param("cart_id")
	itemID := c.Param("item_id")

	if err := h.service.RemoveItem(ctx, cartID, itemID); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}

func (h *Handler) ClearCart(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.clear_cart")
	defer span.End()

	cartID := c.Param("cart_id")
	if err := h.service.ClearCart(ctx, cartID); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

func (h *Handler) MergeCarts(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.merge_carts")
	defer span.End()

	var req struct {
		SourceCartID string `json:"source_cart_id" binding:"required"`
		TargetCartID string `json:"target_cart_id" binding:"required"`
		UserID       string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	if err := h.service.MergeCarts(ctx, req.SourceCartID, req.TargetCartID, req.UserID); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "carts merged"})
}

func (h *Handler) CheckoutPreview(c *gin.Context) {
	ctx, span := otel.Tracer("shopee-cart").Start(c.Request.Context(), "http.checkout_preview")
	defer span.End()

	cartID := c.Param("cart_id")

	var req struct {
		UserID         string `json:"user_id" binding:"required"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	preview, err := h.service.PrepareCheckout(ctx, cartID, req.UserID, req.IdempotencyKey)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, preview)
}

var errorStatusMap = map[error]int{
	domain.ErrCartNotFound:      http.StatusNotFound,
	domain.ErrCartExpired:       http.StatusGone,
	domain.ErrCartFull:          http.StatusConflict,
	domain.ErrItemNotFound:      http.StatusNotFound,
	domain.ErrInvalidQuantity:   http.StatusBadRequest,
	domain.ErrInvalidCartState:  http.StatusConflict,
	domain.ErrDuplicateItem:     http.StatusConflict,
	domain.ErrMergeConflict:     http.StatusConflict,
	domain.ErrSnapshotNotFound:  http.StatusNotFound,
	domain.ErrIdempotencyConflict: http.StatusConflict,
}

func handleError(c *gin.Context, err error) {
	for domainErr, status := range errorStatusMap {
		if err.Error() == domainErr.Error() || len(err.Error()) >= len(domainErr.Error()) && err.Error()[:len(domainErr.Error())] == domainErr.Error() {
			c.AbortWithStatusJSON(status, gin.H{"error_code": domainErr.Error(), "message": err.Error()})
			return
		}
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
