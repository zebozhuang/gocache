package cache

import (
//	"fmt"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Item struct {
	Value      interface{}
	Expiration *time.Time
}

func (item *Item) Expired() bool {
	if item.Expiration == nil {
		return false
	}
	return item.Expiration.Before(time.Now())
}

type cache struct {
	sync.RWMutex
	items             map[string]*Item
	defaultExpiration time.Duration
}

func (c *cache) Set(k string, v interface{}, d time.Duration) {
	c.Lock()
	c.set(k, v, d)
	c.Unlock()
}

func (c *cache) set(k string, v interface{}, d time.Duration) {
	var e time.Time
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d)
	}

	c.items[k] = &Item{
		Value:      v,
		Expiration: &e,
	}
}

func (c *cache) IncrBy(key string, delta int64) error {
	return nil
}

func (c *cache) Incr(key string) error {
	return c.IncrBy(key, 1)
}

func NewCache(defaultExpiration time.Duration) *cache {
	return new(cache)
}
