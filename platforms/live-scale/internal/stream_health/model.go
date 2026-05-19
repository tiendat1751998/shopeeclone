package stream_health

import "time"

type StreamStatus string

const (
	StreamHealthy  StreamStatus = "healthy"
	StreamDegraded StreamStatus = "degraded"
	StreamDown     StreamStatus = "down"
)

type HealthMetric struct {
	Bitrate    int     `json:"bitrate"`
	FrameRate  float64 `json:"frame_rate"`
	LatencyMs  int     `json:"latency_ms"`
	PacketLoss float64 `json:"packet_loss"`
	Viewers    int     `json:"viewers"`
	RecordedAt time.Time `json:"recorded_at"`
}

type AlertThreshold struct {
	MaxLatencyMs    int     `json:"max_latency_ms"`
	MaxPacketLoss   float64 `json:"max_packet_loss"`
	MinBitrate      int     `json:"min_bitrate"`
	MinFrameRate    float64 `json:"min_frame_rate"`
}

type StreamHealth struct {
	StreamID   string         `json:"stream_id"`
	Status     StreamStatus   `json:"status"`
	NodeID     string         `json:"node_id"`
	Region     string         `json:"region"`
	Metrics    []HealthMetric `json:"metrics"`
	Threshold  AlertThreshold `json:"threshold"`
	LastChecked time.Time     `json:"last_checked"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
