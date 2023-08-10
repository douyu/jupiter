package xlru

import (
	"fmt"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/cache/xlru/simplelru"
)

const (
	// Default2QRecentRatio is the ratio of the 2Q cache dedicated
	// to recently added entries that have only been accessed once.
	Default2QRecentRatio = 0.25

	// Default2QGhostEntries is the default ratio of ghost
	// entries kept to track entries recently evicted
	Default2QGhostEntries = 0.50
)

// TwoQueueCache is a thread-safe fixed size 2Q cache.
// 2Q is an enhancement over the standard LRU cache
// in that it tracks both frequently and recently used
// entries separately. This avoids a burst in access to new
// entries from evicting frequently used entries. It adds some
// additional tracking overhead to the standard LRU cache, and is
// computationally about 2x the cost, and adds some metadata over
// head. The ARCCache is similar, but does not require setting any
// parameters.
type TwoQueueCache[K comparable, V any] struct {
	size       int
	recentSize int

	recent      *simplelru.LRU[K, V]
	frequent    *simplelru.LRU[K, V]
	recentEvict *simplelru.LRU[K, struct{}]
	lock        sync.RWMutex
}

// New2Q creates a new TwoQueueCache using the default
// values for the parameters.
func New2Q[K comparable, V any](size int) (*TwoQueueCache[K, V], error) {
	return New2QWithExpire[K, V](size, 0)
}

// New2QParams creates a new TwoQueueCache using the provided
// parameter values.
func New2QParams[K comparable, V any](size int, recentRatio float64, ghostRatio float64) (*TwoQueueCache[K, V], error) {
	return New2QParamsWithExpire[K, V](size, 0, recentRatio, ghostRatio)
}

// New2QWithExpire creates a new TwoQueueCache using the default
// values for the parameters with expire feature.
func New2QWithExpire[K comparable, V any](size int, expire time.Duration) (*TwoQueueCache[K, V], error) {
	return New2QParamsWithExpire[K, V](size, expire, Default2QRecentRatio, Default2QGhostEntries)
}

// New2QParamsWithExpire creates a new TwoQueueCache using the provided
// parameter values with expire feature.
func New2QParamsWithExpire[K comparable, V any](size int, expire time.Duration, recentRatio float64, ghostRatio float64) (*TwoQueueCache[K, V], error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size")
	}
	if recentRatio < 0.0 || recentRatio > 1.0 {
		return nil, fmt.Errorf("invalid recent ratio")
	}
	if ghostRatio < 0.0 || ghostRatio > 1.0 {
		return nil, fmt.Errorf("invalid ghost ratio")
	}

	// Determine the sub-sizes
	recentSize := int(float64(size) * recentRatio)
	evictSize := int(float64(size) * ghostRatio)

	// Allocate the LRUs
	recent, err := simplelru.NewLRUWithExpire[K, V](size, expire, nil)
	if err != nil {
		return nil, err
	}
	frequent, err := simplelru.NewLRUWithExpire[K, V](size, expire, nil)
	if err != nil {
		return nil, err
	}
	recentEvict, err := simplelru.NewLRUWithExpire[K, struct{}](evictSize, expire, nil)
	if err != nil {
		return nil, err
	}

	// Initialize the cache
	c := &TwoQueueCache[K, V]{
		size:        size,
		recentSize:  recentSize,
		recent:      recent,
		frequent:    frequent,
		recentEvict: recentEvict,
	}
	return c, nil
}

func (c *TwoQueueCache[K, V]) Get(key K) (value V, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if this is a frequent value
	if val, ok := c.frequent.Get(key); ok {
		return val, ok
	}

	// If the value is contained in recent, then we
	// promote it to frequent
	if val, expire, ok := c.recent.PeekWithExpireTime(key); ok {
		c.recent.Remove(key)
		var expireDuration time.Duration
		if expire != nil {
			expireDuration = time.Until(*expire)
			if expireDuration < 0 {
				return val, false
			}
		}
		c.frequent.AddEx(key, val, expireDuration)
		return val, ok
	}

	// No hit
	return
}

func (c *TwoQueueCache[K, V]) Add(key K, value V) {
	c.AddEx(key, value, 0)
}

func (c *TwoQueueCache[K, V]) AddEx(key K, value V, expire time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the value is frequently used already,
	// and just update the value
	if c.frequent.Contains(key) {
		c.frequent.AddEx(key, value, expire)
		return
	}

	// Check if the value is recently used, and promote
	// the value into the frequent list
	if c.recent.Contains(key) {
		c.recent.Remove(key)
		c.frequent.AddEx(key, value, expire)
		return
	}

	// If the value was recently evicted, add it to the
	// frequently used list
	if c.recentEvict.Contains(key) {
		c.ensureSpace(true)
		c.recentEvict.Remove(key)
		c.frequent.AddEx(key, value, expire)
		return
	}

	// Add to the recently seen list
	c.ensureSpace(false)
	c.recent.AddEx(key, value, expire)
}

// ensureSpace is used to ensure we have space in the cache
func (c *TwoQueueCache[K, V]) ensureSpace(recentEvict bool) {
	// If we have space, nothing to do
	recentLen := c.recent.Len()
	freqLen := c.frequent.Len()
	if recentLen+freqLen < c.size {
		return
	}

	// If the recent buffer is larger than
	// the target, evict from there
	if recentLen > 0 && (recentLen > c.recentSize || (recentLen == c.recentSize && !recentEvict)) {
		k, _, _ := c.recent.RemoveOldest()
		c.recentEvict.Add(k, struct{}{})
		return
	}

	// Remove from the frequent list otherwise
	c.frequent.RemoveOldest()
}

func (c *TwoQueueCache[K, V]) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.recent.Len() + c.frequent.Len()
}

func (c *TwoQueueCache[K, V]) Keys() []K {
	c.lock.RLock()
	defer c.lock.RUnlock()
	k1 := c.frequent.Keys()
	k2 := c.recent.Keys()
	return append(k1, k2...)
}

func (c *TwoQueueCache[K, V]) Remove(key K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.frequent.Remove(key) {
		return
	}
	if c.recent.Remove(key) {
		return
	}
	if c.recentEvict.Remove(key) {
		return
	}
}

func (c *TwoQueueCache[K, V]) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.recent.Purge()
	c.frequent.Purge()
	c.recentEvict.Purge()
}

func (c *TwoQueueCache[K, V]) Contains(key K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.frequent.Contains(key) || c.recent.Contains(key)
}

func (c *TwoQueueCache[K, V]) Peek(key K) (value V, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, ok := c.frequent.Peek(key); ok {
		return val, ok
	}
	return c.recent.Peek(key)
}
