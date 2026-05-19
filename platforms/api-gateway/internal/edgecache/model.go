package edgecache

import "time"

type CacheEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	TTL       int       `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	HitCount  int64     `json:"hit_count"`
}

type CachePolicy struct {
	PathPattern string   `json:"path_pattern"`
	TTLSeconds  int      `json:"ttl_seconds"`
	VaryBy      VaryBy   `json:"vary_by"`
}

type VaryBy struct {
	Headers     []string `json:"headers"`
	QueryParams []string `json:"query_params"`
}

type CacheStats struct {
	HitCount  int64   `json:"hit_count"`
	MissCount int64   `json:"miss_count"`
	Ratio     float64 `json:"ratio"`
	Entries   int     `json:"entries"`
}
