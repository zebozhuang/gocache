package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = -1
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
	items map[string]*Item
}

func (c *cache) Set(k string, v interface{}, d time.Duration) {
	c.Lock()
	c.set(k, v, d)
	c.Unlock()
}

func (c *cache) set(k string, v interface{}, d time.Duration) {
	var e time.Time
	var item Item

	item.Value = v

	if d >= 0 {
		e = time.Now().Add(d)
		item.Expiration = &e
	}

	c.items[k] = &item
}

func (c *cache) IncrBy(key string, delta int64) (int64, error) {
	c.Lock()
	item, ok := c.items[key]
	if !ok || item.Expired() { // key does not exist or expired.
		c.set(key, delta, NoExpiration)
		c.Unlock()
		return delta, nil
	}

	value, ok := item.Value.(int64)
	if !ok {
		return 0, errors.New(fmt.Sprintf("value type is not integer."))
	}

	newValue := value + delta
	item.Value = newValue
	c.Unlock()
	return newValue, nil
}

func (c *cache) Incr(key string) (int64, error) {
	return c.IncrBy(key, 1)
}

func (c *cache) Get(key string) (interface{}, error) {
	c.RLock()
	v, err := c.get(key)
	c.RUnlock()

	if err != nil {
		return nil, err
	}

	return v.Value, nil
}

func (c *cache) get(key string) (*Item, error) {
	v, ok := c.items[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("key %s does not exist.", key))
	}

	if v.Expired() {
		return nil, errors.New(fmt.Sprintf("key %s has expired.", key))
	}
	return v, nil
}

func (c *cache) Expire(key string, delta time.Duration) error {
	c.Lock()
	item, ok := c.items[key]

	if !ok {
		c.Unlock()
		return errors.New(fmt.Sprintf("key %s does not exist.", key))
	}

	if item.Expired() {
		c.Unlock()
		return errors.New(fmt.Sprintf("key %s has expired.", key))
	}

	e := time.Now().Add(delta)
	item.Expiration = &e

	c.Unlock()
	return nil
}

func (c *cache) ExpireAt(key string, expire *time.Time) error {
	c.Lock()
	item, ok := c.items[key]

	if !ok {
		c.Unlock()
		return errors.New(fmt.Sprintf("key %s does not exist.", key))
	}

	if item.Expired() {
		c.Unlock()
		return errors.New(fmt.Sprintf("key %s has expired.", key))
	}

	item.Expiration = expire
	c.Unlock()
	return nil
}

func (c *cache) Decr(key string) (int64, error) {
	return c.DecrBy(key, 1)
}

func (c *cache) DecrBy(key string, delta int64) (int64, error) {
	c.Lock()
	item, ok := c.items[key]

	if !ok || item.Expired() { // key does not exist or has expired.
		c.set(key, -delta, NoExpiration)
		c.Unlock()
		return -delta, nil
	}

	value, ok := item.Value.(int64)
	if !ok {
		return 0, errors.New(fmt.Sprintf("key %s's value is not integer.", key))
	}

	newValue := value - delta
	item.Value = newValue
	c.Unlock()
	return newValue, nil
}

func NewCache() *cache {
	c := new(cache)
	c.items = map[string]*Item{}

	return c
}
