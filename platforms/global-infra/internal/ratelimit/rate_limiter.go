package ratelimit

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	entries map[string][]time.Time
	mu      sync.RWMutex
}

func NewSlidingWindowCounter() *SlidingWindowCounter {
	return &SlidingWindowCounter{
		entries: make(map[string][]time.Time),
	}
}

func (sw *SlidingWindowCounter) prune(key string, window time.Duration) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	entries, ok := sw.entries[key]
	if !ok {
		return
	}
	cutoff := time.Now().Add(-window)
	var valid []time.Time
	for _, t := range entries {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	sw.entries[key] = valid
}

func (sw *SlidingWindowCounter) count(key string, window time.Duration) int {
	sw.prune(key, window)
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return len(sw.entries[key])
}

func (sw *SlidingWindowCounter) record(key string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.entries[key] = append(sw.entries[key], time.Now())
}

func (sw *SlidingWindowCounter) reset(key string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	delete(sw.entries, key)
}

type RateLimiter struct {
	repo    Repository
	counter *SlidingWindowCounter
}

func NewRateLimiter(repo Repository) *RateLimiter {
	return &RateLimiter{
		repo:    repo,
		counter: NewSlidingWindowCounter(),
	}
}

func (rl *RateLimiter) Check(ctx context.Context, key, strategy string) (*RateLimitCheckResponse, error) {
	rules, err := rl.repo.ListRules(ctx)
	if err != nil {
		return nil, err
	}
	var matchedRule *RateLimitRule
	for _, rule := range rules {
		if ruleMatches(rule, key, Strategy(strategy)) {
			matchedRule = rule
			break
		}
	}
	if matchedRule == nil {
		return &RateLimitCheckResponse{
			Key:       key,
			Allowed:   true,
			Remaining: math.MaxInt32,
			Limit:     math.MaxInt32,
			WindowSec: 60,
		}, nil
	}
	window := time.Duration(matchedRule.WindowSeconds) * time.Second
	currentCount := rl.counter.count(key, window)
	remaining := matchedRule.MaxRequests - currentCount
	if remaining < 0 {
		remaining = 0
	}
	allowed := currentCount < matchedRule.MaxRequests
	return &RateLimitCheckResponse{
		Key:       key,
		Allowed:   allowed,
		Remaining: remaining,
		Limit:     matchedRule.MaxRequests,
		WindowSec: matchedRule.WindowSeconds,
	}, nil
}

func (rl *RateLimiter) Record(ctx context.Context, key, strategy string) {
	rl.counter.record(key)
}

func (rl *RateLimiter) GetRemaining(ctx context.Context, key, strategy string) (int, error) {
	rules, err := rl.repo.ListRules(ctx)
	if err != nil {
		return 0, err
	}
	var matchedRule *RateLimitRule
	for _, rule := range rules {
		if ruleMatches(rule, key, Strategy(strategy)) {
			matchedRule = rule
			break
		}
	}
	if matchedRule == nil {
		return math.MaxInt32, nil
	}
	window := time.Duration(matchedRule.WindowSeconds) * time.Second
	currentCount := rl.counter.count(key, window)
	remaining := matchedRule.MaxRequests - currentCount
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (rl *RateLimiter) Reset(ctx context.Context, key string) {
	rl.counter.reset(key)
}

func (rl *RateLimiter) CreateRule(ctx context.Context, rule *RateLimitRule) error {
	if rule.KeyPattern == "" {
		return fmt.Errorf("key_pattern is required")
	}
	if rule.MaxRequests <= 0 {
		return fmt.Errorf("max_requests must be positive")
	}
	if rule.WindowSeconds <= 0 {
		return fmt.Errorf("window_seconds must be positive")
	}
	if rule.Strategy == "" {
		rule.Strategy = StrategyAPI
	}
	return rl.repo.StoreRule(ctx, rule)
}

func ruleMatches(rule *RateLimitRule, key string, strategy Strategy) bool {
	if rule.Strategy != "" && rule.Strategy != strategy {
		return false
	}
	return patternMatch(rule.KeyPattern, key)
}

func patternMatch(pattern, key string) bool {
	if pattern == "*" || pattern == key {
		return true
	}
	patternLen := len(pattern)
	keyLen := len(key)
	pi, ki := 0, 0
	var starIdx, matchIdx int = -1, 0
	for ki < keyLen {
		if pi < patternLen && (pattern[pi] == key[ki] || pattern[pi] == '?') {
			pi++
			ki++
		} else if pi < patternLen && pattern[pi] == '*' {
			starIdx = pi
			matchIdx = ki
			pi++
		} else if starIdx != -1 {
			pi = starIdx + 1
			matchIdx++
			ki = matchIdx
		} else {
			return false
		}
	}
	for pi < patternLen && pattern[pi] == '*' {
		pi++
	}
	return pi == patternLen
}
