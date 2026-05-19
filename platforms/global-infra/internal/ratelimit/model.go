package ratelimit

import "time"

type Strategy string

const (
	StrategyIP  Strategy = "ip"
	StrategyUser Strategy = "user"
	StrategyAPI Strategy = "api"
)

type RateLimitRule struct {
	KeyPattern    string   `json:"key_pattern"`
	MaxRequests   int      `json:"max_requests"`
	WindowSeconds int      `json:"window_seconds"`
	Strategy      Strategy `json:"strategy"`
}

type SlidingWindowEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

type RateLimitCheckRequest struct {
	Key      string `json:"key" binding:"required"`
	Strategy string `json:"strategy" binding:"required"`
}

type RateLimitCheckResponse struct {
	Key       string `json:"key"`
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	Limit     int    `json:"limit"`
	WindowSec int    `json:"window_seconds"`
}
