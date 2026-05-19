package coordinator

type ShardStatus string

const (
	ShardActive     ShardStatus = "active"
	ShardRelocating ShardStatus = "relocating"
	ShardDead       ShardStatus = "dead"
)

type IndexShard struct {
	ID          string      `json:"id"`
	IndexName   string      `json:"index_name"`
	ShardNumber int         `json:"shard_number"`
	NodeID      string      `json:"node_id"`
	DocCount    int64       `json:"doc_count"`
	SizeBytes   int64       `json:"size_bytes"`
	Status      ShardStatus `json:"status"`
}

type IndexNode struct {
	ID              string  `json:"id"`
	Address         string  `json:"address"`
	Region          string  `json:"region"`
	AvailableShards int     `json:"available_shards"`
	LoadPercentage  float64 `json:"load_percentage"`
	IsActive        bool    `json:"is_active"`
}

type ShardDistribution struct {
	NodeID  string       `json:"node_id"`
	Node    *IndexNode   `json:"node"`
	Shards  []*IndexShard `json:"shards"`
	Count   int          `json:"count"`
}
