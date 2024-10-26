package cache

import (
	"container/list"
	"sync"
	"time"
)

type item[K comparable, V any] struct {
	key    K
	value  V
	expiry time.Time
}

func (i item[K, V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List
	mu       sync.Mutex
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	c := &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		order:    list.New(),
	}

	go func() {
		for range time.Tick(1 * time.Minute) {
			c.mu.Lock()
			for _, element := range c.items {
				cacheItem := element.Value.(*item[K, V])
				if cacheItem.isExpired() {
					c.removeElement(element)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.items[key]; exists {
		cacheItem := element.Value.(*item[K, V])
		cacheItem.value = value
		cacheItem.expiry = time.Now().Add(ttl)
		c.order.MoveToFront(element)
		return
	}

	if c.order.Len() >= c.capacity {
		c.removeOldest()
	}

	cacheItem := &item[K, V]{key: key, value: value, expiry: time.Now().Add(ttl)}
	listElement := c.order.PushFront(cacheItem)
	c.items[key] = listElement

}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, found := c.items[key]
	if !found {
		var zero V
		return zero, false
	}

	cacheItem := element.Value.(*item[K, V])
	if cacheItem.isExpired() {
		c.removeElement(element)
		var zero V
		return zero, false
	}

	c.order.MoveToFront(element)
	return cacheItem.value, true
}

func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		c.removeElement(element)
	}
}

func (c *Cache[K, V]) removeOldest() {
	oldest := c.order.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

func (c *Cache[K, V]) removeElement(element *list.Element) {
	c.order.Remove(element)
	cacheItem := element.Value.(*item[K, V])
	delete(c.items, cacheItem.key)
}
