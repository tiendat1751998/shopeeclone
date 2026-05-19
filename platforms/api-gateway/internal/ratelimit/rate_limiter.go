package ratelimit

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type RateLimiter struct {
	repo        Repository
	entries     map[string][]time.Time
	tokenBucket map[string]int
	lastRefill  map[string]time.Time
	mu          sync.Mutex
}

func NewRateLimiter(repo Repository) *RateLimiter {
	return &RateLimiter{
		repo:        repo,
		entries:     make(map[string][]time.Time),
		tokenBucket: make(map[string]int),
		lastRefill:  make(map[string]time.Time),
	}
}

func (rl *RateLimiter) CreateRule(rule *RateLimitRule) error {
	if rule.Key == "" {
		return fmt.Errorf("key is required")
	}
	if rule.MaxRequests <= 0 {
		return fmt.Errorf("max_requests must be positive")
	}
	if rule.WindowSeconds <= 0 {
		return fmt.Errorf("window_seconds must be positive")
	}
	if rule.BurstSize <= 0 {
		rule.BurstSize = rule.MaxRequests
	}
	return rl.repo.StoreRule(rule)
}

func (rl *RateLimiter) Check(key string) (*CheckResponse, error) {
	rule, err := rl.repo.GetRule(key)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return &CheckResponse{
			Key:       key,
			Allowed:   true,
			Remaining: math.MaxInt32,
			Limit:     math.MaxInt32,
		}, nil
	}

	rl.mu.Lock()
	rl.pruneEntries(key, rule.WindowSeconds)
	currentCount := len(rl.entries[key])
	rl.mu.Unlock()

	tokenAllowed := rl.checkTokenBucket(key, rule.BurstSize, rule.MaxRequests, rule.WindowSeconds)
	windowAllowed := currentCount < rule.MaxRequests
	allowed := windowAllowed && tokenAllowed

	remaining := rule.MaxRequests - currentCount
	if remaining < 0 {
		remaining = 0
	}

	return &CheckResponse{
		Key:       key,
		Allowed:   allowed,
		Remaining: remaining,
		Limit:     rule.MaxRequests,
	}, nil
}

func (rl *RateLimiter) Allow(key string) bool {
	resp, err := rl.Check(key)
	if err != nil {
		return false
	}
	if !resp.Allowed {
		return false
	}
	rl.Record(key)
	return true
}

func (rl *RateLimiter) Record(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.entries[key] = append(rl.entries[key], time.Now())
	rl.consumeToken(key)
}

func (rl *RateLimiter) GetRemaining(key string) (int, error) {
	rule, err := rl.repo.GetRule(key)
	if err != nil {
		return 0, err
	}
	if rule == nil {
		return math.MaxInt32, nil
	}

	rl.mu.Lock()
	rl.pruneEntries(key, rule.WindowSeconds)
	currentCount := len(rl.entries[key])
	rl.mu.Unlock()

	remaining := rule.MaxRequests - currentCount
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.entries, key)
	delete(rl.tokenBucket, key)
	delete(rl.lastRefill, key)
}

func (rl *RateLimiter) pruneEntries(key string, windowSeconds int) {
	cutoff := time.Now().Add(-time.Duration(windowSeconds) * time.Second)
	entries := rl.entries[key]
	var valid []time.Time
	for _, t := range entries {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	rl.entries[key] = valid
}

func (rl *RateLimiter) checkTokenBucket(key string, burstSize, rate, windowSeconds int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	lastRefill, exists := rl.lastRefill[key]
	if !exists {
		rl.tokenBucket[key] = burstSize
		rl.lastRefill[key] = now
		return rl.tokenBucket[key] > 0
	}

	elapsed := now.Sub(lastRefill).Seconds()
	refillRate := float64(burstSize) / float64(windowSeconds)
	newTokens := int(elapsed * refillRate)

	current := rl.tokenBucket[key] + newTokens
	if current > burstSize {
		current = burstSize
	}
	rl.tokenBucket[key] = current
	rl.lastRefill[key] = now

	return current > 0
}

func (rl *RateLimiter) consumeToken(key string) {
	if _, exists := rl.tokenBucket[key]; exists {
		if rl.tokenBucket[key] > 0 {
			rl.tokenBucket[key]--
		}
	}
}
