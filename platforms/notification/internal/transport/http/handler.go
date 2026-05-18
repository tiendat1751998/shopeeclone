package http
import ("net/http"; "strconv"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/notification/internal/application"; "github.com/shopee-clone/shopee/platforms/notification/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.NotificationService }
func NewHandler(s *application.NotificationService) *Handler { return &Handler{service: s} }

func (h *Handler) SendNotification(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-notification").Start(c.Request.Context(), "http.send")
	var req struct { UserID string `json:"user_id" binding:"required"`; Type string `json:"type" binding:"required"`; Title string `json:"title" binding:"required"`; Body string `json:"body" binding:"required"`; Channel string `json:"channel" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	n, err := h.service.SendNotification(ctx, req.UserID, req.Type, req.Title, req.Body, req.Channel, nil)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, n)
}

func (h *Handler) GetInbox(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-notification").Start(c.Request.Context(), "http.inbox")
	userID := c.Query("user_id"); offset := parseInt(c.Query("offset"), 0); limit := parseInt(c.Query("limit"), 20)
	notifs, total, err := h.service.GetInbox(ctx, userID, offset, limit)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"notifications": notifs, "total": total})
}

func (h *Handler) MarkRead(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-notification").Start(c.Request.Context(), "http.mark_read")
	if err := h.service.MarkRead(ctx, c.Param("id"), c.Query("user_id")); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

func (h *Handler) GetUnreadCount(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-notification").Start(c.Request.Context(), "http.unread")
	count, err := h.service.GetUnreadCount(ctx, c.Query("user_id"))
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (h *Handler) UpdatePreference(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-notification").Start(c.Request.Context(), "http.preference")
	var req struct { Channel string `json:"channel" binding:"required"`; Enabled bool `json:"enabled"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	if err := h.service.UpdatePreference(ctx, c.Query("user_id"), req.Channel, req.Enabled); err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, gin.H{"message": "preference updated"})
}

func parseInt(s string, def int) int { if s == "" { return def }; if i, e := strconv.Atoi(s); e == nil { return i }; return def }
var errorStatusMap = map[error]int{domain.ErrNotificationNotFound: http.StatusNotFound, domain.ErrTemplateNotFound: http.StatusNotFound}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() || (len(err.Error()) >= len(e.Error()) && err.Error()[:len(e.Error())] == e.Error()) {
			c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return
		}
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
