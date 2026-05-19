package http

import (
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/collaborative"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/events"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/recommender"
	"github.com/shopee-clone/shopee/platforms/recommendation/internal/trending"
)

type Handler struct {
	recommender recommender.Service
	trending    trending.Service
	collab      collaborative.Service
	publisher   events.Publisher
}

func NewHandler(
	recSvc recommender.Service,
	trendingSvc trending.Service,
	collabSvc collaborative.Service,
	pub events.Publisher,
) *Handler {
	return &Handler{
		recommender: recSvc,
		trending:    trendingSvc,
		collab:      collabSvc,
		publisher:   pub,
	}
}
