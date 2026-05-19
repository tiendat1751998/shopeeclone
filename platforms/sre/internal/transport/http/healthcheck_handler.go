package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateHealthCheckReq struct {
	Name            string `json:"name" binding:"required"`
	Target          string `json:"target" binding:"required"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	Method          string `json:"method"`
	ExpectedStatus  int    `json:"expected_status"`
}

func (h *Handler) CreateHealthCheck(c *gin.Context) {
	var req CreateHealthCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.IntervalSeconds == 0 {
		req.IntervalSeconds = 30
	}
	if req.TimeoutSeconds == 0 {
		req.TimeoutSeconds = 5
	}
	if req.Method == "" {
		req.Method = "GET"
	}
	if req.ExpectedStatus == 0 {
		req.ExpectedStatus = 200
	}
	hc, err := h.healthcheckSvc.CreateCheck(req.Name, req.Target, req.IntervalSeconds, req.TimeoutSeconds, req.Method, req.ExpectedStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, hc)
}

func (h *Handler) RunHealthChecks(c *gin.Context) {
	results, err := h.healthcheckSvc.RunChecks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *Handler) GetHealthCheckResults(c *gin.Context) {
	results, err := h.healthcheckSvc.GetResults()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
