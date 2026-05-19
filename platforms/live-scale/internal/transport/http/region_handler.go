package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetNearestRegion(c *gin.Context) {
	viewerRegion := c.Query("viewer_region")
	if viewerRegion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "viewer_region query parameter required"})
		return
	}
	region, err := h.region.GetNearestRegion(c.Request.Context(), viewerRegion)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, region)
}
