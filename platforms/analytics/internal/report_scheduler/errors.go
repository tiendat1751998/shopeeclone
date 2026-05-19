package report_scheduler

import "errors"

var (
	ErrReportNotFound = errors.New("scheduler: report not found")
	ErrReportInvalid  = errors.New("scheduler: invalid report")
)
