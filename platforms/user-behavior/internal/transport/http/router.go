package http
import ("github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/health"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability")
type Router struct { handler *Handler; health *health.Checker }
func NewRouter(h *Handler, hc *health.Checker) *Router { return &Router{handler: h, health: hc} }
func (r *Router) Setup(e *gin.Engine) {
	e.Use(middleware.Recovery(), middleware.ErrorHandler(), middleware.RequestID(), middleware.CORS(), middleware.OTelMiddleware("shopee-user-behavior"), observability.ObserveHTTPMetrics("shopee-user-behavior"))
	e.GET("/health", r.health.LivenessHandler()); e.GET("/ready", r.health.ReadinessHandler()); e.GET("/metrics", observability.MetricsHandler())
	api := e.Group("/api/v1")
	{ api.POST("/events", r.handler.IngestEvent); api.POST("/events/batch", r.handler.IngestBatch); api.GET("/analytics/trending", r.handler.GetTrending) }
}
