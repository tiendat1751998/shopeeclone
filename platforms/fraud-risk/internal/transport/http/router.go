package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	v1 := engine.Group("/api/v1")
	{
		ruleEngine := v1.Group("/rule-engine")
		{
			ruleEngine.POST("/rules", h.CreateRule)
			ruleEngine.GET("/rules", h.ListRules)
			ruleEngine.POST("/evaluate", h.EvaluateEvent)
			ruleEngine.POST("/rulesets", h.CreateRuleSet)
			ruleEngine.POST("/rulesets/evaluate", h.EvaluateRuleSet)
		}

		riskScore := v1.Group("/risk-score")
		{
			riskScore.POST("/calculate", h.CalculateRiskScore)
			riskScore.GET("/factors", h.GetRiskFactors)
		}

		deviceFp := v1.Group("/device-fingerprint")
		{
			deviceFp.POST("/identify", h.IdentifyDevice)
			deviceFp.POST("/mark-suspicious", h.MarkSuspicious)
		}

		txnMon := v1.Group("/transaction-monitor")
		{
			txnMon.POST("/record", h.RecordTransaction)
			txnMon.POST("/check", h.CheckTransaction)
		}

		behavior := v1.Group("/behavior")
		{
			behavior.POST("/analyze", h.AnalyzeBehavior)
			behavior.POST("/profile", h.BuildProfile)
		}

		v1.POST("/decisions", h.LogDecision)
		v1.GET("/decisions", h.GetDecisions)
	}

	return engine
}
