package replay

import (
	"context"
	"sync"
	"time"
)

type EventBuffer struct {
	mu       sync.RWMutex
	buffers  map[string][]*ReplayEvent
	maxSize  int
	ttl      time.Duration
}

type ReplayEvent struct {
	Sequence  int64       `json:"seq"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"ts"`
}

func NewEventBuffer(maxSize int, ttl time.Duration) *EventBuffer {
	return &EventBuffer{
		buffers: make(map[string][]*ReplayEvent),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

func (eb *EventBuffer) Append(roomID string, seq int64, eventType string, payload interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	entry := &ReplayEvent{
		Sequence:  seq,
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	eb.buffers[roomID] = append(eb.buffers[roomID], entry)
	if len(eb.buffers[roomID]) > eb.maxSize {
		eb.buffers[roomID] = eb.buffers[roomID][len(eb.buffers[roomID])-eb.maxSize:]
	}
}

func (eb *EventBuffer) GetSince(roomID string, sinceSeq int64) []*ReplayEvent {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	buf, ok := eb.buffers[roomID]
	if !ok {
		return nil
	}
	start := 0
	for i, e := range buf {
		if e.Sequence > sinceSeq {
			start = i
			break
		}
	}
	if start >= len(buf) {
		return nil
	}
	result := make([]*ReplayEvent, len(buf)-start)
	copy(result, buf[start:])
	return result
}

func (eb *EventBuffer) Cleanup() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	threshold := time.Now().Add(-eb.ttl)
	for roomID, events := range eb.buffers {
		if len(events) == 0 {
			delete(eb.buffers, roomID)
			continue
		}
		lastEvent := events[len(events)-1]
		if lastEvent.Timestamp.Before(threshold) {
			delete(eb.buffers, roomID)
		}
	}
}

func (eb *EventBuffer) StartCleanup(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				eb.Cleanup()
			}
		}
	}()
}
