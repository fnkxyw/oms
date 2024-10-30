package cache

import (
	"container/list"
	"context"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type item[K comparable, V any] struct {
	key    K
	value  V
	expiry time.Time
	tags   []string
}

func (i item[K, V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type Cache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List
	mu       sync.Mutex
	tagIndex map[string]map[K]struct{}
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	c := &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		order:    list.New(),
		tagIndex: make(map[string]map[K]struct{}),
	}

	go func() {
		for range time.Tick(5 * time.Minute) {
			c.cleanupExpired()
		}
	}()

	return c
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, ttl time.Duration, tags []string) {
	cacheSpan, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Set")
	defer cacheSpan.Finish()

	c.mu.Lock()
	defer c.mu.Unlock()

	cacheSpan.LogKV("action", "Set", "key", key)

	if element, exists := c.items[key]; exists {
		cacheItem := element.Value.(*item[K, V])
		cacheItem.value = value
		cacheItem.expiry = time.Now().Add(ttl)
		cacheItem.tags = tags
		c.order.MoveToFront(element)
		c.updateTagIndex(cacheItem.key, cacheItem.tags)
		cacheSpan.LogKV("info", "updated existing item")
		return
	}

	if c.order.Len() >= c.capacity {
		c.removeOldest()
		cacheSpan.LogKV("info", "removed oldest item to maintain capacity")
	}

	cacheItem := &item[K, V]{key: key, value: value, expiry: time.Now().Add(ttl), tags: tags}
	listElement := c.order.PushFront(cacheItem)
	c.items[key] = listElement
	c.addToTagIndex(cacheItem.key, tags)

	cacheSpan.LogKV("info", "added new item")
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	cacheSpan, ctx := opentracing.StartSpanFromContext(ctx, "Cache.Get")
	defer cacheSpan.Finish()

	cacheSpan.LogKV("action", "Get", "key", key)

	c.mu.Lock()
	defer c.mu.Unlock()

	element, found := c.items[key]
	if !found {
		cacheSpan.LogKV("info", "item not found")
		var zero V
		return zero, false
	}

	cacheItem := element.Value.(*item[K, V])

	if cacheItem.isExpired() {
		c.removeElement(element)
		cacheSpan.LogKV("info", "item expired and removed")
		var zero V
		return zero, false
	}

	c.order.MoveToFront(element)
	cacheSpan.LogKV("info", "item accessed and moved to front")
	return cacheItem.value, true
}

func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.items[key]; found {
		c.removeElement(element)
	}
}

func (c *Cache[K, V]) InvalidateByTags(tags []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, tag := range tags {
		if keys, ok := c.tagIndex[tag]; ok {
			for key := range keys {
				if element, exists := c.items[key]; exists {
					c.removeElement(element)
				}
			}
			delete(c.tagIndex, tag)
		}
	}
}

func (c *Cache[K, V]) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for element := c.order.Back(); element != nil; element = element.Prev() {
		cacheItem := element.Value.(*item[K, V])
		if cacheItem.isExpired() {
			c.removeElement(element)
		}
	}
}

func (c *Cache[K, V]) removeOldest() {
	oldest := c.order.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

func (c *Cache[K, V]) removeElement(element *list.Element) {
	cacheItem := element.Value.(*item[K, V])
	c.order.Remove(element)
	delete(c.items, cacheItem.key)
	c.removeFromTagIndex(cacheItem.key, cacheItem.tags)
}

func (c *Cache[K, V]) addToTagIndex(key K, tags []string) {
	for _, tag := range tags {
		if _, ok := c.tagIndex[tag]; !ok {
			c.tagIndex[tag] = make(map[K]struct{})
		}
		c.tagIndex[tag][key] = struct{}{}
	}
}

func (c *Cache[K, V]) removeFromTagIndex(key K, tags []string) {
	for _, tag := range tags {
		if keys, ok := c.tagIndex[tag]; ok {
			delete(keys, key)
			if len(keys) == 0 {
				delete(c.tagIndex, tag)
			}
		}
	}
}

func (c *Cache[K, V]) updateTagIndex(key K, tags []string) {

	if oldElement, found := c.items[key]; found {
		cacheItem := oldElement.Value.(*item[K, V])
		c.removeFromTagIndex(cacheItem.key, cacheItem.tags)
	}

	c.addToTagIndex(key, tags)
}
