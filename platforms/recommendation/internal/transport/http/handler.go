package http
import ("net/http"; "strconv"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/application"; "github.com/shopee-clone/shopee/platforms/recommendation/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.RecommendationService }
func NewHandler(s *application.RecommendationService) *Handler { return &Handler{service: s} }

func (h *Handler) GetRecommendations(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-recommendation").Start(c.Request.Context(), "http.get_rec")
	userID := c.Query("user_id"); recType := c.Query("type"); limit := parseInt(c.Query("limit"), 20)
	req := domain.RecommendationRequest{UserID: userID, Context: recType, Limit: limit}
	resp, err := h.service.GetRecommendations(ctx, req)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) TrackEvent(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-recommendation").Start(c.Request.Context(), "http.track_event")
	var event domain.UserEvent
	if err := c.ShouldBindJSON(&event); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.TrackEvent(ctx, &event); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "event tracked"})
}

func (h *Handler) GetTrending(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-recommendation").Start(c.Request.Context(), "http.trending")
	limit := parseInt(c.Query("limit"), 20)
	recs, err := h.service.GetTrending(ctx, limit)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"trending": recs})
}

func (h *Handler) GetSimilar(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-recommendation").Start(c.Request.Context(), "http.similar")
	productID := c.Param("product_id"); limit := parseInt(c.Query("limit"), 10)
	recs, err := h.service.GetSimilarProducts(ctx, productID, limit)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"similar": recs})
}

func parseInt(s string, def int) int { if s == "" { return def }; if i, e := strconv.Atoi(s); e == nil { return i }; return def }
var errorStatusMap = map[error]int{domain.ErrRecFailed: http.StatusServiceUnavailable}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() || (len(err.Error()) >= len(e.Error()) && err.Error()[:len(e.Error())] == e.Error()) {
			c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return
		}
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
