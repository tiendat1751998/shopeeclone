package http

import (
	"github.com/shopee-clone/shopee/platforms/search/internal/autocomplete"
	"github.com/shopee-clone/shopee/platforms/search/internal/indexing"
	"github.com/shopee-clone/shopee/platforms/search/internal/query"
	"github.com/shopee-clone/shopee/platforms/search/internal/ranking"
	"github.com/shopee-clone/shopee/platforms/search/internal/search"
)

type Handler struct {
	search       search.Service
	indexing     indexing.Service
	autocomplete autocomplete.Service
	ranking      ranking.Service
	query        query.Service
}

func NewHandler(
	searchSvc search.Service,
	indexingSvc indexing.Service,
	autocompleteSvc autocomplete.Service,
	rankingSvc ranking.Service,
	querySvc query.Service,
) *Handler {
	return &Handler{
		search:       searchSvc,
		indexing:     indexingSvc,
		autocomplete: autocompleteSvc,
		ranking:      rankingSvc,
		query:        querySvc,
	}
}
