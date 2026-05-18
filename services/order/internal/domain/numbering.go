package domain

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OrderNumberGenerator struct {
	mu       sync.Mutex
	sequence int
	lastTime int64
}

func NewOrderNumberGenerator() *OrderNumberGenerator {
	return &OrderNumberGenerator{}
}

func (g *OrderNumberGenerator) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UTC()
	ts := now.UnixMilli()

	if ts == g.lastTime {
		g.sequence++
	} else {
		g.sequence = 0
		g.lastTime = ts
	}

	// Include a random suffix to prevent collisions across restarts
	// Format: ORD-{timestamp}-{sequence}-{random4chars}
	randomSuffix := uuid.New().String()[:4]
	return fmt.Sprintf("ORD-%d-%04d-%s", ts, g.sequence, randomSuffix)
}
