package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type generateEmbeddingRequest struct {
	EntityID   string `json:"entity_id" binding:"required"`
	EntityType string `json:"entity_type" binding:"required"`
	ModelName  string `json:"model_name" binding:"required"`
	Version    string `json:"version" binding:"required"`
}

type findSimilarRequest struct {
	EntityID   string `json:"entity_id" binding:"required"`
	EntityType string `json:"entity_type" binding:"required"`
	TopK       int    `json:"top_k"`
}

func (h *Handler) GenerateEmbedding(c *gin.Context) {
	var req generateEmbeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emb, err := h.embedSvc.GenerateEmbedding(c.Request.Context(), req.EntityID, req.EntityType, req.ModelName, req.Version)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, emb)
}

func (h *Handler) FindSimilar(c *gin.Context) {
	var req findSimilarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topK := req.TopK
	if topK <= 0 {
		topK = 10
	}
	results, err := h.embedSvc.FindSimilar(c.Request.Context(), req.EntityID, req.EntityType, topK)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
