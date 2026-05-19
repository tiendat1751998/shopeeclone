package http

import (
	"github.com/shopee-clone/shopee/platforms/aiml/internal/embeddings"
	"github.com/shopee-clone/shopee/platforms/aiml/internal/experiments"
	"github.com/shopee-clone/shopee/platforms/aiml/internal/featurestore"
	"github.com/shopee-clone/shopee/platforms/aiml/internal/inference"
	"github.com/shopee-clone/shopee/platforms/aiml/internal/modelregistry"
	"github.com/shopee-clone/shopee/platforms/aiml/internal/training"
)

type Handler struct {
	featureSvc   *featurestore.Service
	modelSvc     *modelregistry.Service
	trainingSvc  *training.Service
	inferenceSvc *inference.Service
	embedSvc     *embeddings.Service
	experimentSvc *experiments.Service
}

func NewHandler(
	featureSvc *featurestore.Service,
	modelSvc *modelregistry.Service,
	trainingSvc *training.Service,
	inferenceSvc *inference.Service,
	embedSvc *embeddings.Service,
	experimentSvc *experiments.Service,
) *Handler {
	return &Handler{
		featureSvc:    featureSvc,
		modelSvc:      modelSvc,
		trainingSvc:   trainingSvc,
		inferenceSvc:  inferenceSvc,
		embedSvc:      embedSvc,
		experimentSvc: experimentSvc,
	}
}
