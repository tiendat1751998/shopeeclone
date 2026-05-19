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

func (r *Router) Setup(e *gin.Engine) {
	e.Use(gin.Recovery())
	e.Use(gin.Logger())

	v1 := e.Group("/api/v1")
	{
		v1.POST("/payments", r.handler.ProcessPayment)
		v1.POST("/payments/:id/refund", r.handler.RefundPayment)
		v1.GET("/payments/:id", r.handler.GetPayment)
		v1.GET("/payments", r.handler.ListPayments)
		v1.GET("/payments/order/:orderId", r.handler.GetPaymentsByOrder)

		v1.POST("/ledger/transactions", r.handler.PostLedgerTransaction)
		v1.GET("/ledger/accounts/:id", r.handler.GetLedgerAccount)
		v1.GET("/ledger/accounts/:id/statement", r.handler.GetAccountStatement)
		v1.GET("/ledger/accounts/:id/balance", r.handler.GetAccountBalance)

		v1.POST("/payouts", r.handler.CreatePayout)
		v1.POST("/payouts/batch", r.handler.BatchPayout)
		v1.GET("/payouts/:id", r.handler.GetPayout)
		v1.GET("/payouts", r.handler.ListPayouts)

		v1.POST("/reconciliation/run", r.handler.RunReconciliation)
		v1.GET("/reconciliation/runs", r.handler.ListReconciliationRuns)
		v1.GET("/reconciliation/runs/:id/unmatched", r.handler.GetUnmatchedItems)

		v1.POST("/disputes", r.handler.OpenDispute)
		v1.POST("/disputes/:id/resolve", r.handler.ResolveDispute)
		v1.GET("/disputes", r.handler.ListDisputes)
	}
}
