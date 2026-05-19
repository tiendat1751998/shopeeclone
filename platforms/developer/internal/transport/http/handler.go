package http

import (
	apikeysSvc "github.com/shopee-clone/shopee/platforms/developer/internal/apikeys"
	"github.com/shopee-clone/shopee/platforms/developer/internal/cicd"
	"github.com/shopee-clone/shopee/platforms/developer/internal/docs"
	"github.com/shopee-clone/shopee/platforms/developer/internal/onboarding"
	"github.com/shopee-clone/shopee/platforms/developer/internal/sdk"
	"github.com/shopee-clone/shopee/platforms/developer/internal/webhooks"
)

type Handler struct {
	apikeysSvc    *apikeysSvc.Service
	docsSvc       *docs.Service
	sdkSvc        *sdk.Service
	webhookSvc    *webhooks.Service
	cicdSvc       *cicd.Service
	onboardingSvc *onboarding.Service
}

func NewHandler(
	apikeysSvc *apikeysSvc.Service,
	docsSvc *docs.Service,
	sdkSvc *sdk.Service,
	webhookSvc *webhooks.Service,
	cicdSvc *cicd.Service,
	onboardingSvc *onboarding.Service,
) *Handler {
	return &Handler{
		apikeysSvc:    apikeysSvc,
		docsSvc:       docsSvc,
		sdkSvc:        sdkSvc,
		webhookSvc:    webhookSvc,
		cicdSvc:       cicdSvc,
		onboardingSvc: onboardingSvc,
	}
}
