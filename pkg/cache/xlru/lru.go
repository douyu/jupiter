package xlru

// This package provides a simple LRU cache. It is based on the
// LRU implementation in groupcache:
// https://github.com/golang/groupcache/tree/master/lru

import (
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/cache/xlru/simplelru"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache[K comparable, V any] struct {
	lru  *simplelru.LRU[K, V]
	lock sync.RWMutex
}

// New creates an LRU of the given size
func New[K comparable, V any](size int) (*Cache[K, V], error) {
	return NewWithEvict[K, V](size, nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewWithEvict[K comparable, V any](size int, onEvicted func(key K, value V)) (*Cache[K, V], error) {
	lru, err := simplelru.NewLRU(size, simplelru.EvictCallback[K, V](onEvicted))
	if err != nil {
		return nil, err
	}
	c := &Cache[K, V]{
		lru: lru,
	}
	return c, nil
}

// NewWithExpire constructs a fixed size cache with expire feature
func NewWithExpire[K comparable, V any](size int, expire time.Duration) (*Cache[K, V], error) {
	lru, err := simplelru.NewLRUWithExpire[K, V](size, expire, nil)
	if err != nil {
		return nil, err
	}
	c := &Cache[K, V]{
		lru: lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *Cache[K, V]) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *Cache[K, V]) Add(key K, value V) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Add(key, value)
}

// AddEx adds a value to the cache.  Returns true if an eviction occurred.
func (c *Cache[K, V]) AddEx(key K, value V, expire time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.AddEx(key, value, expire)
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Get(key)
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *Cache[K, V]) Contains(key K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Contains(key)
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache[K, V]) Peek(key K) (V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Peek(key)
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache[K, V]) ContainsOrAdd(key K, value V) (ok, evict bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	} else {
		evict := c.lru.Add(key, value)
		return false, evict
	}
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache[K, V]) PeekOrAdd(key K, value V) (previous V, ok, evict bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, ok := c.lru.Peek(key); ok {
		return v, true, false
	} else {
		var v V
		evict := c.lru.Add(key, value)
		return v, false, evict
	}
}

// Remove removes the provided key from the cache.
func (c *Cache[K, V]) Remove(key K) {
	c.lock.Lock()
	c.lru.Remove(key)
	c.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache[K, V]) RemoveOldest() {
	c.lock.Lock()
	c.lru.RemoveOldest()
	c.lock.Unlock()
}

// Resize changes the cache size.
func (c *Cache[K, V]) Resize(size int) {
	c.lock.Lock()
	c.lru.Resize(size)
	c.lock.Unlock()
}

// GetOldest returns the oldest entry
func (c *Cache[K, V]) GetOldest() (K, V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.GetOldest()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache[K, V]) Keys() []K {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Keys()
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *Cache[K, V]) Values() []V {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Values()
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lru.Len()
}
