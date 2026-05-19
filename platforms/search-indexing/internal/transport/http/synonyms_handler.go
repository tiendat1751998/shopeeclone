package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createSynonymSetRequest struct {
	Words    []string `json:"words"`
	Language string   `json:"language"`
}

func (h *Handler) CreateSynonymSet(c *gin.Context) {
	var req createSynonymSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	set, err := h.synonyms.CreateSet(c.Request.Context(), req.Words, req.Language)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, set)
}

type expandQueryRequest struct {
	Query string `json:"query"`
}

func (h *Handler) ExpandQuery(c *gin.Context) {
	var req expandQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	expanded, err := h.synonyms.ExpandQuery(c.Request.Context(), req.Query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"original": req.Query, "expanded": expanded})
}

func (h *Handler) ListSynonymSets(c *gin.Context) {
	sets, err := h.synonyms.ListSets(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"synonym_sets": sets})
}
