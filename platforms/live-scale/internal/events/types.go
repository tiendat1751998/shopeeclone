package events

import (
	"context"
	"time"
)

type EventType string

const (
	NodeJoined       EventType = "node.joined"
	NodeLeft         EventType = "node.left"
	StreamDegraded   EventType = "stream.degraded"
	StreamDown       EventType = "stream.down"
	CDNPurged        EventType = "cdn.purged"
	RegionFailover   EventType = "region.failover"
	BroadcastMessage EventType = "broadcast.message"
	AlertGenerated   EventType = "alert.generated"
)

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Source    string    `json:"source"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type Producer interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context, handler func(Event) error) error
	Close() error
}
