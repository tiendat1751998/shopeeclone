package funnel

import "errors"

var (
	ErrFunnelInvalid = errors.New("funnel: invalid funnel definition")
	ErrFunnelNotFound = errors.New("funnel: funnel not found")
)
