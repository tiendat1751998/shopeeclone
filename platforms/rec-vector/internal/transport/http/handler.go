package http

import (
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/collabvector"
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/itemembedding"
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/realtime"
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/similarity"
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/userembedding"
	"github.com/shopee-clone/shopee/platforms/rec-vector/internal/vectorstore"
)

type Handler struct {
	vectorStore    vectorstore.VectorStore
	userEmbSvc     *userembedding.Service
	itemEmbSvc     *itemembedding.Service
	similaritySvc  *similarity.Service
	collabSvc      *collabvector.Service
	realtimeSvc    *realtime.Service
}

func NewHandler(
	vs vectorstore.VectorStore,
	ue *userembedding.Service,
	ie *itemembedding.Service,
	sim *similarity.Service,
	collab *collabvector.Service,
	rt *realtime.Service,
) *Handler {
	return &Handler{
		vectorStore:   vs,
		userEmbSvc:    ue,
		itemEmbSvc:    ie,
		similaritySvc: sim,
		collabSvc:     collab,
		realtimeSvc:   rt,
	}
}
