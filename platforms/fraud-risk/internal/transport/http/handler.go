package http

import (
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/behavior"
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/decisionlog"
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/devicefp"
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/riskscoring"
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/ruleengine"
	"github.com/shopee-clone/shopee/platforms/fraud-risk/internal/transactionmon"
)

type Handler struct {
	ruleEngine    *ruleengine.Engine
	riskCalc      *riskscoring.Calculator
	deviceSvc     *devicefp.Service
	txnMon        *transactionmon.Monitor
	behavAnalyzer *behavior.Analyzer
	decLogger     *decisionlog.Logger
}

func NewHandler(
	re *ruleengine.Engine,
	rc *riskscoring.Calculator,
	ds *devicefp.Service,
	tm *transactionmon.Monitor,
	ba *behavior.Analyzer,
	dl *decisionlog.Logger,
) *Handler {
	return &Handler{
		ruleEngine:    re,
		riskCalc:      rc,
		deviceSvc:     ds,
		txnMon:        tm,
		behavAnalyzer: ba,
		decLogger:     dl,
	}
}
