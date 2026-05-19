package edgecache

import (
	"fmt"
	"path"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	mu         sync.RWMutex
	entries    map[string]*CacheEntry
	policies   map[string]*CachePolicy
	hitCount   int64
	missCount  int64
}

func NewCache() *Cache {
	return &Cache{
		entries:  make(map[string]*CacheEntry),
		policies: make(map[string]*CachePolicy),
	}
}

func (c *Cache) Get(key string) (*CacheEntry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		c.missCount++
		return nil, nil
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		c.missCount++
		return nil, nil
	}

	entry.HitCount++
	c.hitCount++
	return entry, nil
}

func (c *Cache) Set(key string, value string, ttl int) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	now := time.Now()
	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		TTL:       ttl,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(ttl) * time.Second),
		HitCount:  0,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry
	return nil
}

func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
	return nil
}

func (c *Cache) PurgeByPattern(pattern string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.entries {
		matched, err := path.Match(pattern, key)
		if err == nil && matched {
			delete(c.entries, key)
			count++
		}
	}
	return count, nil
}

func (c *Cache) GetStats() *CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalRequests := c.hitCount + c.missCount
	var ratio float64
	if totalRequests > 0 {
		ratio = float64(c.hitCount) / float64(totalRequests)
	}

	return &CacheStats{
		HitCount:  c.hitCount,
		MissCount: c.missCount,
		Ratio:     ratio,
		Entries:   len(c.entries),
	}
}

func (c *Cache) AddPolicy(policy *CachePolicy) error {
	if policy.PathPattern == "" {
		return fmt.Errorf("path_pattern is required")
	}
	if policy.TTLSeconds <= 0 {
		return fmt.Errorf("ttl_seconds must be positive")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.policies[policy.PathPattern] = policy
	return nil
}

func (c *Cache) GetPolicy(pattern string) (*CachePolicy, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	policy, ok := c.policies[pattern]
	if !ok {
		return nil, nil
	}
	return policy, nil
}

func (c *Cache) BuildCacheKey(path string, headers map[string]string, queryParams map[string]string, policy *CachePolicy) string {
	var parts []string
	parts = append(parts, path)

	if policy != nil {
		for _, h := range policy.VaryBy.Headers {
			if val, ok := headers[strings.ToLower(h)]; ok {
				parts = append(parts, fmt.Sprintf("h:%s=%s", h, val))
			}
		}
		for _, q := range policy.VaryBy.QueryParams {
			if val, ok := queryParams[q]; ok {
				parts = append(parts, fmt.Sprintf("q:%s=%s", q, val))
			}
		}
	}

	return strings.Join(parts, "|")
}
