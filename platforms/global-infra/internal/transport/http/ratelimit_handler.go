package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type checkRateLimitRequest struct {
	Key      string `json:"key" binding:"required"`
	Strategy string `json:"strategy" binding:"required"`
}

func (h *Handler) CheckRateLimit(c *gin.Context) {
	var req checkRateLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.RateLimiter.Check(c.Request.Context(), req.Key, req.Strategy)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.RateLimiter.Record(c.Request.Context(), req.Key, req.Strategy)

	c.JSON(http.StatusOK, resp)
}
