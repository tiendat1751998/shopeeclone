package recommender

import "errors"

var (
	ErrNoRecommendations = errors.New("no recommendations found")
	ErrInvalidType       = errors.New("invalid recommendation type")
	ErrEmptyInput        = errors.New("empty input: user_id or product_id required")
	ErrServiceUnavailable = errors.New("recommendation service unavailable")
)
