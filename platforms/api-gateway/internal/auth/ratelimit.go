package auth

import (
	"sync"
	"time"
)

type KeyRateLimiter struct {
	mu      sync.Mutex
	records map[string]*keyRecord
}

type keyRecord struct {
	count    int
	windowStart time.Time
}

func NewKeyRateLimiter() *KeyRateLimiter {
	return &KeyRateLimiter{
		records: make(map[string]*keyRecord),
	}
}

func (krl *KeyRateLimiter) RateLimitByKey(key string, maxRequests int, window time.Duration) bool {
	krl.mu.Lock()
	defer krl.mu.Unlock()

	record, exists := krl.records[key]
	now := time.Now()

	if !exists || now.Sub(record.windowStart) >= window {
		krl.records[key] = &keyRecord{
			count:       1,
			windowStart: now,
		}
		return true
	}

	if record.count >= maxRequests {
		return false
	}

	record.count++
	return true
}
