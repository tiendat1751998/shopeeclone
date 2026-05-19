package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type trackEventRequest struct {
	CampaignID string  `json:"campaign_id" binding:"required"`
	UserID     string  `json:"user_id"`
	Revenue    float64 `json:"revenue,omitempty"`
}

func (h *Handler) TrackOpen(c *gin.Context) {
	var req trackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reportingSvc.TrackOpen(c.Request.Context(), req.CampaignID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

func (h *Handler) TrackClick(c *gin.Context) {
	var req trackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reportingSvc.TrackClick(c.Request.Context(), req.CampaignID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

func (h *Handler) GetCampaignReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	report, err := h.reportingSvc.GetCampaignReport(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetAggregatedReport(c *gin.Context) {
	report, err := h.reportingSvc.GetAggregatedReport(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
