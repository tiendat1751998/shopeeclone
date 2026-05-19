package engagement

import (
	"context"
	"sync"
	"time"
	"github.com/shopee-clone/shopee/platforms/live-commerce/internal/fanout"
	redi "github.com/shopee-clone/shopee/platforms/live-commerce/internal/infrastructure/redis"
)

type Counters struct {
	mu       sync.RWMutex
	redis    *redi.Store
	fanout   *fanout.Broadcaster
	rooms    map[string]*RoomCounters
}

type RoomCounters struct {
	Viewers   int64
	Likes     int64
	Gifts     int64
	Shares    int64
	Reactions map[string]int64
	updatedAt time.Time
}

func NewCounters(redis *redi.Store, f *fanout.Broadcaster) *Counters {
	return &Counters{
		redis:  redis,
		fanout: f,
		rooms:  make(map[string]*RoomCounters),
	}
}

func (c *Counters) AddViewer(ctx context.Context, roomID, userID string) {
	c.redis.AddViewer(ctx, roomID, userID)
	if count, err := c.redis.GetViewerCount(ctx, roomID); err == nil {
		c.fanout.Broadcast(ctx, roomID, "viewer_count", map[string]int64{"count": count}, "")
	}
}

func (c *Counters) RemoveViewer(ctx context.Context, roomID, userID string) {
	c.redis.RemoveViewer(ctx, roomID, userID)
	if count, err := c.redis.GetViewerCount(ctx, roomID); err == nil {
		c.fanout.Broadcast(ctx, roomID, "viewer_count", map[string]int64{"count": count}, "")
	}
}

func (c *Counters) AddReaction(ctx context.Context, roomID, reactionType string) {
	c.mu.Lock()
	rc, ok := c.rooms[roomID]
	if !ok {
		rc = &RoomCounters{Reactions: make(map[string]int64), updatedAt: time.Now()}
		c.rooms[roomID] = rc
	}
	rc.Reactions[reactionType]++
	rc.updatedAt = time.Now()
	c.mu.Unlock()
	go func() {
		c.redis.IncrementReaction(ctx, roomID, reactionType)
		c.fanout.Broadcast(ctx, roomID, "reaction", map[string]interface{}{
			"type":  reactionType,
			"count": rc.Reactions[reactionType],
		}, "")
	}()
}

func (c *Counters) AddGift(ctx context.Context, roomID string, amount int64) {
	c.mu.Lock()
	rc, ok := c.rooms[roomID]
	if !ok {
		rc = &RoomCounters{Reactions: make(map[string]int64), updatedAt: time.Now()}
		c.rooms[roomID] = rc
	}
	rc.Gifts += amount
	rc.updatedAt = time.Now()
	c.mu.Unlock()
	go func() {
		c.redis.AddGiftAmount(ctx, roomID, amount)
		c.fanout.Broadcast(ctx, roomID, "gift_total", map[string]int64{"total": rc.Gifts}, "")
	}()
}

func (c *Counters) GetViewerCount(ctx context.Context, roomID string) int64 {
	count, err := c.redis.GetViewerCount(ctx, roomID)
	if err != nil {
		return 0
	}
	return count
}

func (c *Counters) GetReactionCounts(ctx context.Context, roomID string) map[string]int64 {
	c.mu.RLock()
	rc, ok := c.rooms[roomID]
	c.mu.RUnlock()
	if !ok {
		return make(map[string]int64)
	}
	result := make(map[string]int64)
	for k, v := range rc.Reactions {
		result[k] = v
	}
	return result
}

func (c *Counters) GetGiftTotal(ctx context.Context, roomID string) int64 {
	c.mu.RLock()
	rc, ok := c.rooms[roomID]
	c.mu.RUnlock()
	if !ok {
		return 0
	}
	return rc.Gifts
}
