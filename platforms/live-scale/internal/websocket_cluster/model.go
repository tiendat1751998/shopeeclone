package websocket_cluster

import "time"

type WSNodeStatus string

const (
	WSNodeActive   WSNodeStatus = "active"
	WSNodeDegraded WSNodeStatus = "degraded"
	WSNodeDown     WSNodeStatus = "down"
)

type WSNode struct {
	ID           string       `json:"id"`
	Address      string       `json:"address"`
	Region       string       `json:"region"`
	Status       WSNodeStatus `json:"status"`
	RoomCount    int          `json:"room_count"`
	ClientCount  int          `json:"client_count"`
	MaxRooms     int          `json:"max_rooms"`
	MaxClients   int          `json:"max_clients"`
	LastHeartbeat time.Time   `json:"last_heartbeat"`
	RegisteredAt  time.Time   `json:"registered_at"`
}

type WSCluster struct {
	Nodes map[string]*WSNode `json:"nodes"`
	Rooms map[string]string  `json:"rooms"`
}

type RoomAssignment struct {
	RoomID   string `json:"room_id"`
	NodeID   string `json:"node_id"`
	Clients  int    `json:"clients"`
	AssignedAt time.Time `json:"assigned_at"`
}

type ClientRouting struct {
	ClientID string `json:"client_id"`
	RoomID   string `json:"room_id"`
	NodeID   string `json:"node_id"`
}
