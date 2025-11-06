package shared

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type CacheItem struct {
	Value     any
	ExpiresAt time.Time
}

type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]CacheItem),
		ttl:  ttl,
	}

	go c.cleanupRoutine(1 * time.Minute)

	return c
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	log.Info().Str("key", key).Msg("cache set")
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		log.Info().Str("key", key).Msg("cache miss")
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		log.Info().Str("key", key).Msg("cache expired")
		c.Clear(key)
		return nil, false
	}

	log.Info().Str("key", key).Msg("cache hit")
	return item.Value, true
}

func (c *Cache) Clear(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) cleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		now := time.Now()

		c.mu.Lock()
		for key, item := range c.data {
			if now.After(item.ExpiresAt) {
				delete(c.data, key)
				log.Debug().Str("key", key).Msg("cache key cleaned up (expired)")
			}
		}
		c.mu.Unlock()
	}
}
