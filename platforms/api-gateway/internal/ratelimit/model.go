package ratelimit

type RateLimitRule struct {
	Key           string `json:"key"`
	MaxRequests   int    `json:"max_requests"`
	WindowSeconds int    `json:"window_seconds"`
	BurstSize     int    `json:"burst_size"`
}

type CheckResponse struct {
	Key       string `json:"key"`
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	Limit     int    `json:"limit"`
}
