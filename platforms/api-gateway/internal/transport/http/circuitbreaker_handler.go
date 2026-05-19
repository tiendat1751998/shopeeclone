package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCircuitBreakerRequest struct {
	Name                string `json:"name" binding:"required"`
	ServiceName         string `json:"service_name" binding:"required"`
	FailureThreshold    int    `json:"failure_threshold"`
	RecoveryTimeout     int    `json:"recovery_timeout_seconds"`
	HalfOpenMaxRequests int    `json:"half_open_max_requests"`
}

func (h *Handler) CreateCircuitBreaker(c *gin.Context) {
	var req createCircuitBreakerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cb, err := h.CBSvc.Create(req.Name, req.ServiceName, req.FailureThreshold, req.RecoveryTimeout, req.HalfOpenMaxRequests)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cb)
}

func (h *Handler) ListCircuitBreakers(c *gin.Context) {
	cbs, err := h.CBSvc.List()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"circuit_breakers": cbs})
}

func (h *Handler) RecordCircuitBreakerSuccess(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.CBSvc.RecordSuccess(id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success recorded"})
}

func (h *Handler) RecordCircuitBreakerFailure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.CBSvc.RecordFailure(id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "failure recorded"})
}
