package http
import ("net/http"; "strconv"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/advertising/internal/application"; "github.com/shopee-clone/shopee/platforms/advertising/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.AdvertisingService }
func NewHandler(s *application.AdvertisingService) *Handler { return &Handler{service: s} }

func (h *Handler) CreateCampaign(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-advertising").Start(c.Request.Context(), "http.create_campaign")
	var req struct { AdvertiserID string `json:"advertiser_id" binding:"required"`; Name string `json:"name" binding:"required"`; Budget int64 `json:"budget" binding:"required"`; DailyBudget int64 `json:"daily_budget" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	campaign, err := h.service.CreateCampaign(ctx, req.AdvertiserID, req.Name, req.Budget, req.DailyBudget, nil, nil)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, campaign)
}

func (h *Handler) ServeAds(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-advertising").Start(c.Request.Context(), "http.serve_ads")
	query := c.Query("q"); userID := c.Query("user_id"); limit := parseInt(c.Query("limit"), 10)
	ads, err := h.service.ServeAds(ctx, query, userID, limit)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"ads": ads})
}

func (h *Handler) RecordImpression(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-advertising").Start(c.Request.Context(), "http.impression")
	var imp domain.Impression
	if err := c.ShouldBindJSON(&imp); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.RecordImpression(ctx, &imp); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "recorded"})
}

func (h *Handler) RecordClick(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-advertising").Start(c.Request.Context(), "http.click")
	var click domain.Click
	if err := c.ShouldBindJSON(&click); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.RecordClick(ctx, &click); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "recorded"})
}

func parseInt(s string, def int) int { if s == "" { return def }; if i, e := strconv.Atoi(s); e == nil { return i }; return def }
var errorStatusMap = map[error]int{domain.ErrCampaignNotFound: http.StatusNotFound}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() { c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return }
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
