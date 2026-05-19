package monitoring

import "time"

type HealthStatus string

const (
	HealthGreen  HealthStatus = "green"
	HealthYellow HealthStatus = "yellow"
	HealthRed    HealthStatus = "red"
)

type IndexMetric struct {
	IndexName       string       `json:"index_name"`
	ShardCount      int          `json:"shard_count"`
	DocCount        int64        `json:"doc_count"`
	SizeBytes       int64        `json:"size_bytes"`
	IndexingRate    float64      `json:"indexing_rate"`
	SearchLatencyP99 float64     `json:"search_latency_p99"`
	Health          HealthStatus `json:"health"`
	LastUpdated     time.Time    `json:"last_updated"`
}

type ClusterHealth struct {
	NodesCount       int `json:"nodes_count"`
	ActiveShards     int `json:"active_shards"`
	RelocatingShards int `json:"relocating_shards"`
	UnassignedShards int `json:"unassigned_shards"`
}
