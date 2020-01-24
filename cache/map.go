package cache

import (
	"log"
	"sync"
	"time"
)

// MapCache - simplest cache implementation based on map and mutex
type MapCache struct {
	store      map[uint64]*Item
	lock       sync.RWMutex
	expiration time.Duration
}

// Get - deliver cached item if it exists
func (c *MapCache) Get(uri string) (*Item, bool) {
	key := constructKey(uri)
	c.lock.RLock()
	item, ok := c.store[key]
	c.lock.RUnlock()
	return item, ok
}

// Set - set cached item
func (c *MapCache) Set(uri string, item *Item) {
	key := constructKey(uri)
	c.lock.Lock()
	c.store[key] = item
	c.lock.Unlock()
}

// GarbageCollect - each 30 seconds clear stale items,
// should be started as goroutine
func (c *MapCache) GarbageCollect(timeToExpire time.Duration) {
	for {
		select {
		case <-time.After(30 * time.Second):
			c.lock.Lock()
			for key, item := range c.store {
				if item.isExpired(c.expiration) {
					delete(c.store, key)
				}
			}
			c.lock.Unlock()
			log.Println("Garbage Collection")
		}
	}
}
