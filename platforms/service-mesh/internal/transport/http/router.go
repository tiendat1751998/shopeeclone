package http

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *Handler
}

func NewRouter(h *Handler) *Router {
	return &Router{handler: h}
}

func (r *Router) Setup(engine *gin.Engine) {
	api := engine.Group("/api/v1")

	discovery := api.Group("/discovery")
	discovery.POST("/register", r.handler.RegisterService)
	discovery.POST("/heartbeat", r.handler.Heartbeat)
	discovery.GET("/services", r.handler.ListServices)
	discovery.GET("/discover", r.handler.DiscoverServices)

	mtls := api.Group("/mtls")
	mtls.POST("/ca", r.handler.CreateCA)
	mtls.POST("/certificates/issue", r.handler.IssueCertificate)
	mtls.POST("/certificates/renew", r.handler.RenewCertificate)
	mtls.POST("/certificates/revoke", r.handler.RevokeCertificate)
	mtls.GET("/certificates", r.handler.ListCertificates)
	mtls.POST("/certificates/verify", r.handler.VerifyCertificate)

	traffic := api.Group("/traffic")
	traffic.POST("/rules", r.handler.CreateTrafficRule)
	traffic.GET("/rules", r.handler.ListTrafficRules)
	traffic.POST("/evaluate", r.handler.EvaluateRoute)

	lb := api.Group("/loadbalancer")
	lb.POST("/next", r.handler.GetNextInstance)

	resilience := api.Group("/resilience")
	resilience.POST("/retry", r.handler.ExecuteWithRetry)
	resilience.POST("/bulkhead", r.handler.ExecuteWithBulkhead)
	resilience.GET("/circuit-breakers", r.handler.ListCircuitBreakers)

	telemetry := api.Group("/telemetry")
	telemetry.POST("/record-call", r.handler.RecordCall)
	telemetry.GET("/service-graph", r.handler.GetServiceGraph)
	telemetry.GET("/traces", r.handler.GetTraces)
}
