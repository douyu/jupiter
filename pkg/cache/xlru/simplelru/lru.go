package simplelru

import (
	"errors"
	"time"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

// LRU implements a non-thread safe fixed size LRU cache
type LRU[K comparable, V any] struct {
	size      int
	evictList *LruList[K, V]
	freeList  *LruList[K, V]
	items     map[K]*Entry[K, V]
	expire    time.Duration
	onEvict   EvictCallback[K, V]
}

// NewLRU constructs an LRU of the given size
func NewLRU[K comparable, V any](size int, onEvict EvictCallback[K, V]) (*LRU[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	c := &LRU[K, V]{
		size:      size,
		evictList: NewList[K, V](),
		freeList:  NewList[K, V](),
		items:     make(map[K]*Entry[K, V]),
		expire:    0,
		onEvict:   onEvict,
	}
	for i := 0; i < size; i++ {
		var (
			k K
			v V
		)
		c.freeList.PushFront(k, v)
	}
	return c, nil
}

// NewLRUWithExpire contrusts an LRU of the given size and expire time
func NewLRUWithExpire[K comparable, V any](size int, expire time.Duration, onEvict EvictCallback[K, V]) (*LRU[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	c := &LRU[K, V]{
		size:      size,
		evictList: NewList[K, V](),
		freeList:  NewList[K, V](),
		items:     make(map[K]*Entry[K, V]),
		expire:    expire,
		onEvict:   onEvict,
	}
	for i := 0; i < size; i++ {
		var (
			k K
			v V
		)
		c.freeList.PushFront(k, v)
	}

	return c, nil
}

// Purge is used to completely clear the cache
func (c *LRU[K, V]) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
	c.freeList.Init()
	for i := 0; i < c.size; i++ {
		var (
			k K
			v V
		)
		c.freeList.PushFront(k, v)
	}
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU[K, V]) Add(key K, value V) bool {
	return c.AddEx(key, value, 0)
}

// AddEx adds a value to the cache with expire.  Returns true if an eviction occurred.
func (c *LRU[K, V]) AddEx(key K, value V, expire time.Duration) bool {
	var ex *time.Time = nil
	if expire > 0 {
		expire := time.Now().Add(expire)
		ex = &expire
	} else if c.expire > 0 {
		expire := time.Now().Add(c.expire)
		ex = &expire
	}
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value = value
		ent.ExpiresAt = ex
		return false
	}

	evict := c.evictList.Length() >= c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}

	// Add new item
	ent := c.freeList.Front()
	ent.Key = key
	ent.Value = value
	ent.ExpiresAt = ex
	c.freeList.Remove(ent)
	c.evictList.PushElementFront(ent)
	c.items[key] = ent

	return evict
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	if ent, ok := c.items[key]; ok {
		if ent.IsExpired() {
			var v V
			return v, false
		}
		c.evictList.MoveToFront(ent)
		return ent.Value, true
	}
	return
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	if ent, ok := c.items[key]; ok {
		if ent.IsExpired() {
			return false
		}
		return ok
	}
	return
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	v, _, ok := c.PeekWithExpireTime(key)
	return v, ok
}

// Returns the key value (or undefined if not found) and its associated expire
// time without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) PeekWithExpireTime(key K) (
	value V, expire *time.Time, ok bool) {
	var v V
	if ent, ok := c.items[key]; ok {
		if ent.IsExpired() {
			return v, nil, false
		}
		return ent.Value, ent.ExpiresAt, true
	}
	return v, nil, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) bool {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value, true
	}
	return key, value, false
}

// GetOldest returns the oldest entry
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		return ent.Key, ent.Value, true
	}
	return key, value, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU[K, V]) Keys() []K {
	keys := make([]K, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		keys[i] = ent.Key
		i++
	}
	return keys
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *LRU[K, V]) Values() []V {
	values := make([]V, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		values[i] = ent.Value
		i++
	}
	return values
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	return c.evictList.Length()
}

// Resize changes the cache size.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

// removeOldest removes the oldest item from the cache.
func (c *LRU[K, V]) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRU[K, V]) removeElement(e *Entry[K, V]) {
	c.evictList.Remove(e)
	c.freeList.PushElementFront(e)
	delete(c.items, e.Key)
	if c.onEvict != nil {
		c.onEvict(e.Key, e.Value)
	}
}
