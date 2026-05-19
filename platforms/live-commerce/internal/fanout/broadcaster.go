package fanout

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"github.com/shopee-clone/shopee/packages/go-shared/pkg/observability"
	"github.com/shopee-clone/shopee/platforms/live-commerce/internal/websocket"
	"go.uber.org/zap"
)

type BroadcastMsg struct {
	RoomID    string
	Event     string
	Payload   interface{}
	ExcludeID string
}

type Broadcaster struct {
	hub    *websocket.Hub
	input  chan *BroadcastMsg
	mu     sync.RWMutex
	global map[string]chan struct{}
}

func NewBroadcaster(hub *websocket.Hub) *Broadcaster {
	b := &Broadcaster{
		hub:    hub,
		input:  make(chan *BroadcastMsg, 1024),
		global: make(map[string]chan struct{}),
	}
	go b.process()
	return b
}

func (b *Broadcaster) process() {
	for msg := range b.input {
		wrapper := map[string]interface{}{
			"type":      msg.Event,
			"payload":   msg.Payload,
			"timestamp": time.Now().UnixMilli(),
		}
		data, err := json.Marshal(wrapper)
		if err != nil {
			observability.LogWithTrace(nil).Error("broadcast marshal", zap.Error(err))
			continue
		}
		b.hub.BroadcastToRoomBytes(msg.RoomID, data, msg.ExcludeID)
	}
}

func (b *Broadcaster) Broadcast(ctx context.Context, roomID, event string, payload interface{}, excludeID string) {
	select {
	case b.input <- &BroadcastMsg{RoomID: roomID, Event: event, Payload: payload, ExcludeID: excludeID}:
	default:
		observability.LogWithTrace(ctx).Warn("broadcast buffer full",
			zap.String("room", roomID), zap.String("event", event))
	}
}

func (b *Broadcaster) BroadcastSync(roomID, event string, payload interface{}, excludeID string) {
	wrapper := map[string]interface{}{
		"type":      event,
		"payload":   payload,
		"timestamp": time.Now().UnixMilli(),
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return
	}
	b.hub.BroadcastToRoomBytes(roomID, data, excludeID)
}

func (b *Broadcaster) RoomCount(roomID string) int {
	return b.hub.RoomCount(roomID)
}

func (b *Broadcaster) Stop() {
	close(b.input)
}
