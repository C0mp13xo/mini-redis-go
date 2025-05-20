package main

import (
	"container/list"
	"encoding/json"
	"os"
	"time"
)

type Snapshot struct {
	Items []CacheItem `json:"items"`
}

func (c *Cache) SaveSnapshot(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := Snapshot{Items: make([]CacheItem, 0, len(c.store))}
	for _, elem := range c.store {
		item := elem.Value.(*CacheItem)
		snapshot.Items = append(snapshot.Items, *item)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(snapshot)
}

func (c *Cache) LoadSnapshot(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	var snapshot Snapshot
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&snapshot)
	if err != nil {
		return err
	}

	c.store = make(map[string]*list.Element)
	c.order.Init()
	for _, item := range snapshot.Items {
		if time.Now().After(item.ExpiresAt) {
			continue
		}
		element := c.order.PushFront(&item)
		c.store[item.Key] = element
	}
	return nil
}
