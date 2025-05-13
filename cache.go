package main

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	Value  string
	Expiry time.Time
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]Entry
}

func NewCache() *Cache {
	c := &Cache{
		store: make(map[string]Entry),
	}
	return c
}

// Background cleanup for Expired Keys
func (c *Cache) cleanupExpiredKeys(done <-chan struct{}) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			c.mu.Lock()
			currTime := time.Now()
			for k, v := range c.store {
				if !v.Expiry.IsZero() && currTime.After(v.Expiry) {
					delete(c.store, k)
				}
			}
			c.mu.Unlock()

		case <-done:
			fmt.Println("stopping process to cleanupExpired keys ")
			return
		}

	}

}

func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := Entry{
		Value: value,
	}
	if ttl > 0 {
		entry.Expiry = time.Now().Add(ttl)
	}
	c.store[key] = entry
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return "", false
	}
	if time.Now().After(entry.Expiry) {
		return "", false
	}

	return entry.Value, true
}

func (c *Cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.store[key]; exists {
		delete(c.store, key)
		return true
	}
	return false
}
