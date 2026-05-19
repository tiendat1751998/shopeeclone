package http

import (
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/bulkindexer"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/coordinator"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/monitoring"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/pipeline"
	"github.com/shopee-clone/shopee/platforms/search-indexing/internal/synonyms"
)

type Handler struct {
	coordinator coordinator.Service
	bulkindexer bulkindexer.Service
	pipeline    pipeline.Service
	synonyms    synonyms.Service
	monitoring  monitoring.Service
}

func NewHandler(
	coordinator coordinator.Service,
	bulkindexer bulkindexer.Service,
	pipeline pipeline.Service,
	synonyms synonyms.Service,
	monitoring monitoring.Service,
) *Handler {
	return &Handler{
		coordinator: coordinator,
		bulkindexer: bulkindexer,
		pipeline:    pipeline,
		synonyms:    synonyms,
		monitoring:  monitoring,
	}
}
