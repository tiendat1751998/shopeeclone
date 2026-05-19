package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/events"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/metrics"
)

type feedbackRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	EventType string `json:"event_type"`
	SessionID string `json:"session_id"`
}

func (h *Handler) RecordFeedback(c *gin.Context) {
	var req feedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.UserID == "" || req.ProductID == "" || req.EventType == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user_id, product_id, and event_type are required"})
		return
	}

	metrics.EventsTracked.Inc()

	switch req.EventType {
	case "click":
		h.collab.RecordInteraction(context.Background(), req.UserID, req.ProductID, false)
		h.trending.RecordInteraction(context.Background(), req.ProductID)
		h.publisher.Publish(context.Background(), events.EventItemClicked, events.ItemClicked{
			UserID:    req.UserID,
			ProductID: req.ProductID,
			SessionID: req.SessionID,
			Timestamp: time.Now(),
		})

	case "purchase":
		h.collab.RecordInteraction(context.Background(), req.UserID, req.ProductID, true)
		h.trending.RecordInteraction(context.Background(), req.ProductID)
		h.publisher.Publish(context.Background(), events.EventItemPurchased, events.ItemPurchased{
			UserID:    req.UserID,
			ProductID: req.ProductID,
			Timestamp: time.Now(),
		})

	case "view":
		h.trending.RecordInteraction(context.Background(), req.ProductID)

	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unknown event_type: must be click, purchase, or view"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feedback recorded"})
}
