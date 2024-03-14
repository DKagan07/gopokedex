package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheVal struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu sync.Mutex
	// the key in our cacheEntry is the URL
	// the createdAt time is now(), and the val is the raw []bytes
	cacheEntry map[string]cacheVal
	duration   time.Duration
}

func NewCache(interval time.Duration) Cache {
	m := make(map[string]cacheVal)

	c := Cache{
		cacheEntry: m,
		duration:   interval,
	}

	// start reaping
	c.reapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cacheEntry[key]; ok {
		fmt.Println("already in cache")
	}

	value := cacheVal{
		time.Now(),
		val,
	}

	c.cacheEntry[key] = value
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.cacheEntry[key]
	if !ok {
		return nil, false
	}

	return v.val, true
}

func (c *Cache) reapLoop() {
	// now := time.Now()
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			<-ticker.C
			// logic here
			c.mu.Lock()
			for k, v := range c.cacheEntry {
				ttk := v.createdAt.Add(c.duration)
				if time.Now().After(ttk) {
					fmt.Println("deleted entry: ", k)
					delete(c.cacheEntry, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
