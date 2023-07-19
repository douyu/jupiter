package xfreecache

type cache[K comparable, V any] struct {
	storage[K, V]
}

// GetAndSetCacheData 获取缓存后数据
func (c *cache[K, V]) GetAndSetCacheData(key string, id K, fn func() (V, error)) (value V, err error) {
	resMap, err := c.GetAndSetCacheMap(key, []K{id}, func([]K) (map[K]V, error) {
		innerVal, innerErr := fn()
		return map[K]V{id: innerVal}, innerErr
	})
	value = resMap[id]
	return
}

// GetCacheValue 获取缓存数据
func (c *cache[K, V]) GetCacheValue(key string, id K) (value V) {
	resMap, _ := c.getCacheMap(key, []K{id})
	value = resMap[id]
	return
}

// SetCacheValue 设置缓存数据
func (c *cache[K, V]) SetCacheValue(key string, id K, fn func() (V, error)) (err error) {
	err = c.setCacheMap(key, []K{id}, func([]K) (map[K]V, error) {
		innerVal, innerErr := fn()
		return map[K]V{id: innerVal}, innerErr
	}, nil)
	return
}

// GetAndSetCacheMap 获取缓存后数据 map形式
func (c *cache[K, V]) GetAndSetCacheMap(key string, ids []K, fn func([]K) (map[K]V, error)) (v map[K]V, err error) {
	// 获取缓存数据
	v, idsNone := c.getCacheMap(key, ids)

	// 设置缓存数据
	err = c.setCacheMap(key, idsNone, fn, v)
	return
}

// GetCacheMap 获取缓存数据 map形式
func (c *cache[K, V]) GetCacheMap(key string, ids []K) (v map[K]V) {
	v, _ = c.getCacheMap(key, ids)
	return
}

// SetCacheMap 设置缓存数据 map形式
func (c *cache[K, V]) SetCacheMap(key string, ids []K, fn func([]K) (map[K]V, error)) (err error) {
	err = c.setCacheMap(key, ids, fn, nil)
	return
}
