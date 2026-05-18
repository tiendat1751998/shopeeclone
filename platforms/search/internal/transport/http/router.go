package http
import ("github.com/gin-gonic/gin"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/health"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/middleware"; "github.com/shopee-clone/shopee/packages/go-shared/pkg/observability")
type Router struct { handler *Handler; health *health.Checker }
func NewRouter(h *Handler, hc *health.Checker) *Router { return &Router{handler: h, health: hc} }
func (r *Router) Setup(e *gin.Engine) {
	e.Use(middleware.Recovery(), middleware.ErrorHandler(), middleware.RequestID(), middleware.CORS(), middleware.OTelMiddleware("shopee-search"), observability.ObserveHTTPMetrics("shopee-search"))
	e.GET("/health", r.health.LivenessHandler()); e.GET("/ready", r.health.ReadinessHandler()); e.GET("/metrics", observability.MetricsHandler())
	api := e.Group("/api/v1")
	{ api.GET("/search", r.handler.Search); api.GET("/autocomplete", r.handler.Autocomplete); api.GET("/trending", r.handler.GetTrending); api.POST("/index", r.handler.IndexProduct) }
}
