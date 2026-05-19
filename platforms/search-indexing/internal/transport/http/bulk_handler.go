package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createBulkJobRequest struct {
	IndexName      string `json:"index_name"`
	TotalDocuments int    `json:"total_documents"`
}

func (h *Handler) CreateBulkJob(c *gin.Context) {
	var req createBulkJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	job, err := h.bulkindexer.CreateJob(c.Request.Context(), req.IndexName, req.TotalDocuments)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

type submitBatchRequest struct {
	JobID       string                   `json:"job_id"`
	Documents   []map[string]interface{} `json:"documents"`
	BatchNumber int                      `json:"batch_number"`
}

func (h *Handler) SubmitBatch(c *gin.Context) {
	var req submitBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	batch, err := h.bulkindexer.SubmitBatch(c.Request.Context(), req.JobID, req.Documents, req.BatchNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, batch)
}

func (h *Handler) ListBulkJobs(c *gin.Context) {
	jobs, err := h.bulkindexer.ListJobs(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func (h *Handler) GetBulkJobProgress(c *gin.Context) {
	id := c.Param("id")
	job, err := h.bulkindexer.GetJobProgress(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}
