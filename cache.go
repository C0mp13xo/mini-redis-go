package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// type Entry struct {
// 	Value  string
// 	Expiry time.Time
// }

type CacheItem struct {
	Key       string
	Value     string
	ExpiresAt time.Time
}

type Cache struct {
	mu       sync.RWMutex
	store    map[string]*list.Element
	order    *list.List
	capacity int
}

func NewCache(capacity int) *Cache {
	c := &Cache{
		store:    make(map[string]*list.Element),
		order:    list.New(),
		capacity: capacity,
	}
	return c
}

func (c *Cache) evictLRU() {
	lruElement := c.order.Back()
	if lruElement != nil {
		c.removeElement(lruElement)
	}
}

// Background cleanup for Expired Keys
func (c *Cache) cleanupExpiredKeys(interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("cleaning up expired keys")
			c.removeExpiredKeys()
		case <-stopChan:
			return
		}
	}
}

func (c *Cache) removeExpiredKeys() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for e := c.order.Back(); e != nil; {
		prev := e.Prev()
		item := e.Value.(*CacheItem)
		if time.Now().After(item.ExpiresAt) {
			c.removeElement(e)
		}
		e = prev
	}

}

func (c *Cache) removeElement(e *list.Element) {
	item := e.Value.(*CacheItem)
	delete(c.store, item.Key)
	c.order.Remove(e)
}

func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.store[key]; ok {
		item := elem.Value.(*CacheItem)
		item.Value = value
		item.ExpiresAt = time.Now().Add(ttl)
		c.order.MoveToFront(elem)
		return
	}

	if len(c.store) >= c.capacity {
		c.evictLRU()
	}

	item := &CacheItem{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	element := c.order.PushFront(item)
	c.store[key] = element
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.store[key]
	if !ok {
		return "", false
	}
	item := elem.Value.(*CacheItem)
	if time.Now().After(item.ExpiresAt) {
		//key has expired, remove from map
		c.removeElement(elem)
		return "", false
	}
	c.order.MoveToFront(elem)
	return item.Value, true
}

func (c *Cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.store[key]; !ok {

		return false
	} else {
		c.removeElement(elem)
		return true
	}
}
