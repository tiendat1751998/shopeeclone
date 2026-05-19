package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type setCacheRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"`
}

type purgeCacheRequest struct {
	Pattern string `json:"pattern" binding:"required"`
}

func (h *Handler) GetCacheValue(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	entry, err := h.EdgeCache.Get(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entry == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "cache miss"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *Handler) SetCacheValue(c *gin.Context) {
	var req setCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TTL <= 0 {
		req.TTL = 300
	}

	if err := h.EdgeCache.Set(req.Key, req.Value, req.TTL); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cached"})
}

func (h *Handler) PurgeCache(c *gin.Context) {
	var req purgeCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.EdgeCache.PurgeByPattern(req.Pattern)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"purged": count})
}

func (h *Handler) GetCacheStats(c *gin.Context) {
	stats := h.EdgeCache.GetStats()
	c.JSON(http.StatusOK, stats)
}
