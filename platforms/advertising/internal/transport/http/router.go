package http
import ("github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/health"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability")
type Router struct { handler *Handler; health *health.Checker }
func NewRouter(h *Handler, hc *health.Checker) *Router { return &Router{handler: h, health: hc} }
func (r *Router) Setup(e *gin.Engine) {
	e.Use(middleware.Recovery(), middleware.ErrorHandler(), middleware.RequestID(), middleware.CORS(), middleware.OTelMiddleware("shopee-advertising"), observability.ObserveHTTPMetrics("shopee-advertising"))
	e.GET("/health", r.health.LivenessHandler()); e.GET("/ready", r.health.ReadinessHandler()); e.GET("/metrics", observability.MetricsHandler())
	api := e.Group("/api/v1")
	{ api.POST("/ads/campaigns", r.handler.CreateCampaign); api.GET("/ads/serve", r.handler.ServeAds); api.POST("/ads/impressions", r.handler.RecordImpression); api.POST("/ads/clicks", r.handler.RecordClick) }
}
