package region

import "time"

type RegionStatus string

const (
	RegionActive  RegionStatus = "active"
	RegionDegraded RegionStatus = "degraded"
	RegionDown    RegionStatus = "down"
)

type Region struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Code      string       `json:"code"`
	Status    RegionStatus `json:"status"`
	Latency   int          `json:"latency_ms"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type LatencyMap struct {
	FromRegion string `json:"from_region"`
	ToRegion   string `json:"to_region"`
	LatencyMs  int    `json:"latency_ms"`
}

type GeoRoutingRule struct {
	ID            string       `json:"id"`
	ViewerRegion  string       `json:"viewer_region"`
	PrimaryRegion string       `json:"primary_region"`
	FailoverRegion string      `json:"failover_region"`
	Priority      int          `json:"priority"`
	IsActive      bool         `json:"is_active"`
	CreatedAt     time.Time    `json:"created_at"`
}
