package http

import (
	"github.com/shopee-clone/shopee/platforms/notification/internal/dispatcher"
	"github.com/shopee-clone/shopee/platforms/notification/internal/email"
	"github.com/shopee-clone/shopee/platforms/notification/internal/inapp"
	"github.com/shopee-clone/shopee/platforms/notification/internal/notifier"
	"github.com/shopee-clone/shopee/platforms/notification/internal/preferences"
	"github.com/shopee-clone/shopee/platforms/notification/internal/push"
	"github.com/shopee-clone/shopee/platforms/notification/internal/sms"
	"github.com/shopee-clone/shopee/platforms/notification/internal/template"
)

type Handler struct {
	notifier    notifier.Service
	push        push.Service
	email       email.Service
	sms         sms.Service
	inapp       inapp.Service
	preferences preferences.Service
	template    template.Service
	dispatcher  dispatcher.Service
}

func NewHandler(
	notifSvc notifier.Service,
	pushSvc push.Service,
	emailSvc email.Service,
	smsSvc sms.Service,
	inappSvc inapp.Service,
	prefSvc preferences.Service,
	tmplSvc template.Service,
	dispatchSvc dispatcher.Service,
) *Handler {
	return &Handler{
		notifier:    notifSvc,
		push:        pushSvc,
		email:       emailSvc,
		sms:         smsSvc,
		inapp:       inappSvc,
		preferences: prefSvc,
		template:    tmplSvc,
		dispatcher:  dispatchSvc,
	}
}
