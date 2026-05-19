package deliveryopt

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	SavePattern(ctx context.Context, p *UserEngagementPattern) error
	GetPattern(ctx context.Context, userID string, channel string) (*UserEngagementPattern, error)
	ListPatterns(ctx context.Context) ([]*UserEngagementPattern, error)
	Enqueue(ctx context.Context, msg *QueuedMessage) error
	Dequeue(ctx context.Context) (*QueuedMessage, error)
	ListQueue(ctx context.Context) ([]*QueuedMessage, error)
	RecordSend(ctx context.Context, channel string) error
	GetChannelSendCount(ctx context.Context, channel string, since time.Time) (int, error)
}

type InMemoryRepository struct {
	mu        sync.RWMutex
	patterns  map[string]*UserEngagementPattern
	queue     []*QueuedMessage
	sendCount map[string][]time.Time
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		patterns:  make(map[string]*UserEngagementPattern),
		queue:     make([]*QueuedMessage, 0),
		sendCount: make(map[string][]time.Time),
	}
}

func (r *InMemoryRepository) SavePattern(ctx context.Context, p *UserEngagementPattern) error {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := p.UserID + ":" + p.Channel
	r.patterns[key] = p
	return nil
}

func (r *InMemoryRepository) GetPattern(ctx context.Context, userID string, channel string) (*UserEngagementPattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := userID + ":" + channel
	p, ok := r.patterns[key]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *InMemoryRepository) ListPatterns(ctx context.Context) ([]*UserEngagementPattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*UserEngagementPattern
	for _, p := range r.patterns {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryRepository) Enqueue(ctx context.Context, msg *QueuedMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	msg.CreatedAt = time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queue = append(r.queue, msg)
	return nil
}

func (r *InMemoryRepository) Dequeue(ctx context.Context) (*QueuedMessage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.queue) == 0 {
		return nil, nil
	}

	highestIdx := 0
	for i := 1; i < len(r.queue); i++ {
		if r.queue[i].Priority > r.queue[highestIdx].Priority {
			highestIdx = i
		}
	}

	msg := r.queue[highestIdx]
	r.queue = append(r.queue[:highestIdx], r.queue[highestIdx+1:]...)
	return msg, nil
}

func (r *InMemoryRepository) ListQueue(ctx context.Context) ([]*QueuedMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*QueuedMessage, len(r.queue))
	copy(result, r.queue)
	return result, nil
}

func (r *InMemoryRepository) RecordSend(ctx context.Context, channel string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sendCount[channel] = append(r.sendCount[channel], time.Now())
	return nil
}

func (r *InMemoryRepository) GetChannelSendCount(ctx context.Context, channel string, since time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, t := range r.sendCount[channel] {
		if t.After(since) {
			count++
		}
	}
	return count, nil
}
