package http

import (
	analyticsSvc "github.com/shopee-clone/shopee/platforms/analytics/internal/analytics"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/cohort"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/dashboard"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/events"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/funnel"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/report_scheduler"
	"github.com/shopee-clone/shopee/platforms/analytics/internal/session"
)

type Handler struct {
	analyticsSvc *analyticsSvc.Service
	eventSvc     *events.Service
	funnelSvc    *funnel.Service
	cohortSvc    *cohort.Service
	sessionSvc   *session.Service
	dashboardSvc *dashboard.Service
	scheduleSvc  *report_scheduler.Service
	publisher    events.Publisher
}

func NewHandler(
	analyticsSvc *analyticsSvc.Service,
	eventSvc *events.Service,
	funnelSvc *funnel.Service,
	cohortSvc *cohort.Service,
	sessionSvc *session.Service,
	dashboardSvc *dashboard.Service,
	scheduleSvc *report_scheduler.Service,
	pub events.Publisher,
) *Handler {
	return &Handler{
		analyticsSvc: analyticsSvc,
		eventSvc:     eventSvc,
		funnelSvc:    funnelSvc,
		cohortSvc:    cohortSvc,
		sessionSvc:   sessionSvc,
		dashboardSvc: dashboardSvc,
		scheduleSvc:  scheduleSvc,
		publisher:    pub,
	}
}
