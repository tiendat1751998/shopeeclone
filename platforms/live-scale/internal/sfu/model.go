package sfu

import "time"

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"
	NodeStatusDegraded NodeStatus = "degraded"
	NodeStatusDown     NodeStatus = "down"
)

type SFUNode struct {
	ID          string     `json:"id"`
	Address     string     `json:"address"`
	Region      string     `json:"region"`
	Capacity    int        `json:"capacity"`
	CurrentLoad int        `json:"current_load"`
	StreamCount int        `json:"stream_count"`
	Status      NodeStatus `json:"status"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	RegisteredAt  time.Time `json:"registered_at"`
}

type SFUCluster struct {
	Nodes map[string]*SFUNode `json:"nodes"`
}

type StreamSession struct {
	ID         string    `json:"id"`
	StreamID   string    `json:"stream_id"`
	NodeID     string    `json:"node_id"`
	Region     string    `json:"region"`
	Viewers    int       `json:"viewers"`
	StartedAt  time.Time `json:"started_at"`
	LastActive time.Time `json:"last_active"`
}
