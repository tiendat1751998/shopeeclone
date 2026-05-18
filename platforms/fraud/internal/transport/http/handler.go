package http
import ("net/http"; "github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/platforms/fraud/internal/application"; "github.com/shopee-clone/shopee/platforms/fraud/internal/domain"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"; "go.opentelemetry.io/otel"; "go.uber.org/zap")

type Handler struct { service *application.FraudService }
func NewHandler(s *application.FraudService) *Handler { return &Handler{service: s} }

func (h *Handler) ScoreTransaction(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-fraud").Start(c.Request.Context(), "http.score")
	var req struct { UserID string `json:"user_id" binding:"required"`; OrderID string `json:"order_id" binding:"required"`; Amount int64 `json:"amount" binding:"required"`; DeviceIP string `json:"device_ip"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	score, err := h.service.ScoreTransaction(ctx, req.UserID, req.OrderID, req.Amount, req.DeviceIP)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusOK, score)
}

func (h *Handler) CreateCase(c *gin.Context) {
	ctx, _ := otel.Tracer("shopee-fraud").Start(c.Request.Context(), "http.create_case")
	var req struct { UserID string `json:"user_id" binding:"required"`; OrderID string `json:"order_id" binding:"required"`; Score float64 `json:"score" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_REQUEST", "message": err.Error()}); return }
	case_, err := h.service.CreateCase(ctx, req.UserID, req.OrderID, req.Score)
	if err != nil { handleError(c, err); return }
	c.JSON(http.StatusCreated, case_)
}

var errorStatusMap = map[error]int{domain.ErrFraudDetection: http.StatusServiceUnavailable}
func handleError(c *gin.Context, err error) {
	for e, s := range errorStatusMap {
		if err.Error() == e.Error() { c.AbortWithStatusJSON(s, gin.H{"error_code": e.Error(), "message": err.Error()}); return }
	}
	observability.LogWithTrace(c.Request.Context()).Error("unhandled error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_ERROR", "message": "An unexpected error occurred"})
}
