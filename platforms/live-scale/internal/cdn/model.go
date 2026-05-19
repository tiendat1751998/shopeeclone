package cdn

import "time"

type CDNEndpoint struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Provider    string   `json:"provider"`
	Region      string   `json:"region"`
	LatencyMs   int      `json:"latency_ms"`
	Capacity    int      `json:"capacity"`
	CurrentLoad int      `json:"current_load"`
	Status      string   `json:"status"`
	Populations []string `json:"populations"`
	CreatedAt   time.Time `json:"created_at"`
}

type CDNPurgeRequest struct {
	ID        string   `json:"id"`
	URLs      []string `json:"urls,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Reason    string   `json:"reason"`
	RequestedAt time.Time `json:"requested_at"`
}

type EdgeCacheStatus struct {
	EndpointID    string    `json:"endpoint_id"`
	URL           string    `json:"url"`
	Cached        bool      `json:"cached"`
	TTL           int       `json:"ttl_seconds"`
	LastRefreshed time.Time `json:"last_refreshed"`
	HitCount      int64     `json:"hit_count"`
}
