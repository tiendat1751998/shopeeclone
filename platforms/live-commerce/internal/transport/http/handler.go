package http
import ("net/http"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/live-commerce/internal/application"; "github.com/shopee-clone/shopee/platforms/live-commerce/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.LiveCommerceService }
func NewHandler(s *application.LiveCommerceService) *Handler { return &Handler{service: s} }

func (h *Handler) CreateLivestream(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.create")
	var req struct { SellerID string `json:"seller_id" binding:"required"`; Title string `json:"title" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	ls, err := h.service.CreateLivestream(ctx, req.SellerID, req.Title)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, ls)
}

func (h *Handler) StartLivestream(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.start")
	if err := h.service.StartLivestream(ctx, c.Param("id")); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "livestream started"})
}

func (h *Handler) EndLivestream(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.end")
	if err := h.service.EndLivestream(ctx, c.Param("id")); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "livestream ended"})
}

func (h *Handler) SendMessage(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.chat")
	var req struct { UserID string `json:"user_id" binding:"required"`; Content string `json:"content" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	msg, err := h.service.SendChatMessage(ctx, c.Param("room_id"), req.UserID, req.Content)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, msg)
}

func (h *Handler) SendReaction(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.reaction")
	var req struct { UserID string `json:"user_id" binding:"required"`; Type string `json:"type" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.SendReaction(ctx, c.Param("room_id"), req.UserID, req.Type); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "reaction sent"})
}

func (h *Handler) GetViewerCount(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-live-commerce").Start(c.Request.Context(), "http.viewers")
	count, err := h.service.GetViewerCount(ctx, c.Param("id"))
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"viewer_count": count})
}

var errorStatusMap = map[error]int{domain.ErrLivestreamNotFound: http.StatusNotFound}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() { c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return }
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
