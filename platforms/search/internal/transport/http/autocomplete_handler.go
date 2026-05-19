package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Autocomplete(c *gin.Context) {
	prefix := c.Query("q")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	if prefix == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	result, err := h.autocomplete.Suggest(c.Request.Context(), prefix, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Trending(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	trending, err := h.autocomplete.GetTrending(c.Request.Context(), limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trending": trending})
}

func (h *Handler) RecordSearch(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Query == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	if err := h.autocomplete.RecordSearch(c.Request.Context(), req.Query); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recorded"})
}
