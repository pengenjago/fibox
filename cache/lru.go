package cache

import (
	"context"
	"time"

	"github.com/pengenjago/fibox/logging"

	lru "github.com/hashicorp/golang-lru/v2"
)

// Cache interface defines the operations for our cache wrapper
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}) error
	SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	Clear(ctx context.Context) error
	Stats() Stats
}

// Stats represents cache statistics
type Stats struct {
	Hits   int64
	Misses int64
	Size   int
}

// LRUCache implements the Cache interface using golang-lru
type LRUCache struct {
	cache  *lru.Cache[string, cacheItem]
	stats  Stats
	ttlMap map[string]time.Time
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewLRUCache creates a new LRU cache with the specified size
func NewLRUCache(size int) Cache {
	cache, err := lru.New[string, cacheItem](size)
	if err != nil {
		return nil
	}

	return &LRUCache{
		cache:  cache,
		ttlMap: make(map[string]time.Time),
	}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(ctx context.Context, key string) (interface{}, bool) {
	item, ok := c.cache.Get(key)
	if !ok {
		c.stats.Misses++
		logging.DebugWithFields("Cache miss",
			map[string]interface{}{
				"key":       key,
				"cache_hit": false,
			})
		return nil, false
	}

	// Check if the item has expired
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.cache.Remove(key)
		delete(c.ttlMap, key)
		c.stats.Misses++
		logging.DebugWithFields("Cache expired",
			map[string]interface{}{
				"key":       key,
				"cache_hit": false,
			})
		return nil, false
	}

	c.stats.Hits++
	logging.DebugWithFields("Cache hit",
		map[string]interface{}{
			"key":       key,
			"cache_hit": true,
		})
	return item.value, true
}

// Set stores a value in the cache without TTL
func (c *LRUCache) Set(ctx context.Context, key string, value interface{}) error {
	item := cacheItem{
		value:     value,
		expiresAt: time.Time{}, // Zero time means no expiration
	}
	c.cache.Add(key, item)
	delete(c.ttlMap, key) // Remove any existing TTL for this key

	logging.DebugWithFields("Cache set",
		map[string]interface{}{
			"key": key,
		})
	return nil
}

// SetWithTTL stores a value in the cache with a TTL
func (c *LRUCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	item := cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.cache.Add(key, item)
	c.ttlMap[key] = item.expiresAt

	logging.DebugWithFields("Cache set with TTL",
		map[string]interface{}{
			"key":      key,
			"duration": ttl.String(),
		})
	return nil
}

// Delete removes a value from the cache
func (c *LRUCache) Delete(ctx context.Context, key string) error {
	c.cache.Remove(key)
	delete(c.ttlMap, key)

	logging.DebugWithFields("Cache delete",
		map[string]interface{}{
			"key": key,
		})
	return nil
}

// Clear removes all values from the cache
func (c *LRUCache) Clear(ctx context.Context) error {
	c.cache.Purge()
	c.ttlMap = make(map[string]time.Time)

	logging.DebugWithFields("Cache cleared",
		map[string]interface{}{
			"size": c.cache.Len(),
		})
	return nil
}

// DeleteByPattern removes all cache entries that match the given pattern
func (c *LRUCache) DeleteByPattern(ctx context.Context, pattern string) error {
	keysToDelete := []string{}

	// Get all keys in the cache
	for key := range c.ttlMap {
		// Simple pattern matching - in a real implementation, you might want to use regex
		if c.matchesPattern(key, pattern) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	// Delete matching keys
	for _, key := range keysToDelete {
		c.cache.Remove(key)
		delete(c.ttlMap, key)
	}

	logging.DebugWithFields("Cache delete by pattern",
		map[string]interface{}{
			"pattern": pattern,
			"count":   len(keysToDelete),
		})

	return nil
}

// matchesPattern checks if a key matches a simple pattern
// Supports wildcard '*' at the end of the pattern
func (c *LRUCache) matchesPattern(key, pattern string) bool {
	// Simple implementation for patterns like "location:list:*"
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return key == pattern
}

// Stats returns cache statistics
func (c *LRUCache) Stats() Stats {
	c.stats.Size = c.cache.Len()
	return c.stats
}
