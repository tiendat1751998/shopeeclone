package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/errors"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/domain"
	"github.com/shopee-clone/shopee/services/catalog-product/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Handler struct {
	productUseCase  *usecase.ProductUseCase
	categoryUseCase *usecase.CategoryUseCase
}

func NewHandler(productUC *usecase.ProductUseCase, categoryUC *usecase.CategoryUseCase) *Handler {
	return &Handler{
		productUseCase:  productUC,
		categoryUseCase: categoryUC,
	}
}

func (h *Handler) GetProduct(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.product.get")
	defer span.End()

	spuID := c.Param("spu_id")
	if spuID == "" {
		c.Error(domain.ErrInvalidProductData)
		return
	}

	product, err := h.productUseCase.GetByID(ctx, spuID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) ListProducts(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.product.list")
	defer span.End()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	filter := domain.ProductFilter{
		Page:       page,
		Size:       size,
		CategoryID: c.Query("category_id"),
		SellerID:   c.Query("seller_id"),
		Search:     c.Query("search"),
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		SortBy:     c.Query("sort_by"),
	}

	result, err := h.productUseCase.List(ctx, filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": result.Products,
		"total":    result.Total,
		"page":     result.Page,
		"size":     result.Size,
	})
}

func (h *Handler) CreateProduct(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.product.create")
	defer span.End()

	var req struct {
		Title       string              `json:"title" binding:"required"`
		Description string              `json:"description"`
		CategoryID  string              `json:"category_id" binding:"required"`
		SKUs        []domain.SKU        `json:"skus" binding:"required"`
		Attributes  map[string]string   `json:"attributes"`
		Images      []string            `json:"images"`
		SellerID    string              `json:"seller_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewValidation("Invalid request body").WithDetail("body", err.Error()))
		return
	}

	product := &domain.Product{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		SKUs:        req.SKUs,
		Attributes:  req.Attributes,
		Images:      req.Images,
		SellerID:    req.SellerID,
	}

	created, err := h.productUseCase.Create(ctx, product)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.product.update")
	defer span.End()

	spuID := c.Param("spu_id")
	if spuID == "" {
		c.Error(domain.ErrInvalidProductData)
		return
	}

	var req struct {
		Title       string              `json:"title"`
		Description string              `json:"description"`
		CategoryID  string              `json:"category_id"`
		Attributes  map[string]string   `json:"attributes"`
		Images      []string            `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewValidation("Invalid request body"))
		return
	}

	product := &domain.Product{
		SPUID:       spuID,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Attributes:  req.Attributes,
		Images:      req.Images,
	}

	if err := h.productUseCase.Update(ctx, product); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.product.delete")
	defer span.End()

	spuID := c.Param("spu_id")
	if spuID == "" {
		c.Error(domain.ErrInvalidProductData)
		return
	}

	if err := h.productUseCase.Delete(ctx, spuID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *Handler) ListCategories(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.category.list")
	defer span.End()

	parentID := c.Query("parent_id")
	levelStr := c.Query("level")
	var level int32
	if levelStr != "" {
		l, err := strconv.Atoi(levelStr)
		if err == nil {
			level = int32(l)
		}
	}

	categories, err := h.categoryUseCase.List(ctx, parentID, level)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func (h *Handler) GetCategory(c *gin.Context) {
	ctx, span := otel.Tracer("catalog-product").Start(c.Request.Context(), "handler.category.get")
	defer span.End()

	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.Error(domain.ErrInvalidCategory)
		return
	}

	category, err := h.categoryUseCase.GetByID(ctx, categoryID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}
