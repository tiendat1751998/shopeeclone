package multiregion

import "time"

type LatencyConfig struct {
	P95LatencyMs float64 `json:"p95_latency_ms"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
}

type Region struct {
	Name           string                  `json:"name"`
	Code           string                  `json:"code"`
	IsActive       bool                    `json:"is_active"`
	FailoverRegion string                  `json:"failover_region,omitempty"`
	Endpoints      map[string]string       `json:"endpoints"`
	LatencyConfig  *LatencyConfig          `json:"latency_config,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type FailoverResult struct {
	PrimaryRegion   string `json:"primary_region"`
	FailoverRegion  string `json:"failover_region"`
	Endpoint        string `json:"endpoint"`
	IsFailover      bool   `json:"is_failover"`
}
