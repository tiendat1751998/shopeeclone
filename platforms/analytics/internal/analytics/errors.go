package analytics

import "errors"

var (
	ErrQueryInvalid     = errors.New("analytics: invalid query")
	ErrMetricNotFound   = errors.New("analytics: metric not found")
	ErrDimensionInvalid = errors.New("analytics: invalid dimension")
	ErrNoData           = errors.New("analytics: no data available")
	ErrTimeRangeInvalid = errors.New("analytics: invalid time range")
)
