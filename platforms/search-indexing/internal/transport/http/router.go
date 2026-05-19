package http

import "github.com/gin-gonic/gin"

type Router struct {
	handler *Handler
}

func NewRouter(h *Handler) *Router {
	return &Router{handler: h}
}

func (r *Router) Setup(e *gin.Engine) {
	api := e.Group("/api/v1")
	{
		coord := api.Group("/coordinator")
		{
			coord.POST("/nodes", r.handler.RegisterNode)
			coord.GET("/nodes", r.handler.ListNodes)
			coord.POST("/assign-shard", r.handler.AssignShard)
			coord.POST("/rebalance", r.handler.Rebalance)
		}

		bulk := api.Group("/bulk")
		{
			bulk.POST("/jobs", r.handler.CreateBulkJob)
			bulk.POST("/submit", r.handler.SubmitBatch)
			bulk.GET("/jobs", r.handler.ListBulkJobs)
			bulk.GET("/jobs/:id", r.handler.GetBulkJobProgress)
		}

		pl := api.Group("/pipelines")
		{
			pl.POST("", r.handler.CreatePipeline)
			pl.POST("/process", r.handler.ProcessDocument)
			pl.GET("", r.handler.ListPipelines)
		}

		syn := api.Group("/synonyms")
		{
			syn.POST("", r.handler.CreateSynonymSet)
			syn.POST("/expand", r.handler.ExpandQuery)
			syn.GET("", r.handler.ListSynonymSets)
		}

		mon := api.Group("/monitoring")
		{
			mon.POST("/metrics", r.handler.ReportMetrics)
			mon.GET("/cluster-health", r.handler.GetClusterHealth)
			mon.GET("/indexes", r.handler.GetIndexMetrics)
		}
	}
}
