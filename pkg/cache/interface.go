package cache

// Storage 接入本地缓存库需要实现的接口
type Storage[K comparable, V any] interface {
	GetCacheMapOrigin(key string, ids []K) (v map[K]V, idsNone []K)
	SetCacheMapOrigin(key string, idsNone []K, fn func([]K) (map[K]V, error), v map[K]V) (err error)
}

type Cache[K comparable, V any] struct {
	Storage[K, V]
}

// GetAndSetCacheData 获取缓存后数据
func (c *Cache[K, V]) GetAndSetCacheData(key string, id K, fn func() (V, error)) (value V, err error) {
	resMap, err := c.GetAndSetCacheMap(key, []K{id}, func([]K) (map[K]V, error) {
		innerVal, innerErr := fn()
		return map[K]V{id: innerVal}, innerErr
	})
	value = resMap[id]
	return
}

// GetCacheValue 获取缓存数据
func (c *Cache[K, V]) GetCacheValue(key string, id K) (value V) {
	resMap, _ := c.GetCacheMapOrigin(key, []K{id})
	value = resMap[id]
	return
}

// SetCacheValue 设置缓存数据
func (c *Cache[K, V]) SetCacheValue(key string, id K, fn func() (V, error)) (err error) {
	err = c.SetCacheMapOrigin(key, []K{id}, func([]K) (map[K]V, error) {
		innerVal, innerErr := fn()
		return map[K]V{id: innerVal}, innerErr
	}, nil)
	return
}

// GetAndSetCacheMap 获取缓存后数据 map形式
func (c *Cache[K, V]) GetAndSetCacheMap(key string, ids []K, fn func([]K) (map[K]V, error)) (v map[K]V, err error) {
	// 获取缓存数据
	v, idsNone := c.GetCacheMapOrigin(key, ids)

	// 设置缓存数据
	err = c.SetCacheMapOrigin(key, idsNone, fn, v)
	return
}

// GetCacheMap 获取缓存数据 map形式
func (c *Cache[K, V]) GetCacheMap(key string, ids []K) (v map[K]V) {
	v, _ = c.GetCacheMapOrigin(key, ids)
	return
}

// SetCacheMap 设置缓存数据 map形式
func (c *Cache[K, V]) SetCacheMap(key string, ids []K, fn func([]K) (map[K]V, error)) (err error) {
	err = c.SetCacheMapOrigin(key, ids, fn, nil)
	return
}
