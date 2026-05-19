package dashboard

import "errors"

var (
	ErrDashboardNotFound = errors.New("dashboard: not found")
	ErrWidgetNotFound    = errors.New("dashboard: widget not found")
)
