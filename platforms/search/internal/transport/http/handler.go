package http
import ("net/http"; "strconv"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/search/internal/application"; "github.com/shopee-clone/shopee/platforms/search/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.SearchService }
func NewHandler(s *application.SearchService) *Handler { return &Handler{service: s} }

func (h *Handler) Search(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-search").Start(c.Request.Context(), "http.search")
	query := domain.SearchQuery{
		Query: c.Query("q"), CategoryID: c.Query("category_id"), ShopID: c.Query("shop_id"),
		SortBy: c.Query("sort_by"), Page: parseInt(c.Query("page"), 1), Limit: parseInt(c.Query("limit"), 20),
	}
	result, err := h.service.Search(ctx, query)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Autocomplete(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-search").Start(c.Request.Context(), "http.autocomplete")
	prefix := c.Query("q"); limit := parseInt(c.Query("limit"), 10)
	result, err := h.service.Autocomplete(ctx, prefix, limit)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetTrending(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-search").Start(c.Request.Context(), "http.trending")
	queries, err := h.service.GetTrendingQueries(ctx, 20)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"trending": queries})
}

func (h *Handler) IndexProduct(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-search").Start(c.Request.Context(), "http.index")
	var doc domain.IndexDocument
	if err := c.ShouldBindJSON(&doc); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.IndexProduct(ctx, &doc); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "indexed"})
}

func parseInt(s string, def int) int { if s == "" { return def }; if i, e := strconv.Atoi(s); e == nil { return i }; return def }

var errorStatusMap = map[error]int{domain.ErrSearchFailed: http.StatusServiceUnavailable, domain.ErrIndexNotFound: http.StatusNotFound}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() || (len(err.Error()) >= len(e.Error()) && err.Error()[:len(e.Error())] == e.Error()) {
			c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return
		}
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
