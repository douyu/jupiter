package xfreecache

import (
	"github.com/coocood/freecache"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	glru "github.com/hnlq715/golang-lru"
	"go.uber.org/zap"
	"sync"
)

const (
	cacheTypeFreeCache = "freecache"
	cacheTypeLruCache  = "lru"
)

var (
	innerCache    *freecache.Cache
	innerLruCache *glru.ARCCache
	once          sync.Once
	lruOnce       sync.Once
)

type storage[K comparable, V any] interface {
	getCacheMap(key string, ids []K) (v map[K]V, idsNone []K)
	setCacheMap(key string, idsNone []K, fn func([]K) (map[K]V, error), v map[K]V) (err error)
}

func newFreeCache[K comparable, V any](c *Config) (localCache *LocalCache[K, V]) {
	// 初始化缓存实例 如果没有配置大小，默认200000
	lruOnce.Do(func() {
		var err error
		size := cfg.GetInt("jupiter.cache.sizeLru")
		if size <= 0 {
			size = 200000
		}
		innerLruCache, err = glru.NewARCWithExpire(size, c.Expire)
		if err != nil {
			xlog.Jupiter().Panic("glru NewARCWithExpire error", xlog.FieldErr(err))
		}
	})
	return &LocalCache[K, V]{
		cache: cache[K, V]{&myLru[K, V]{
			config: c,
			cache:  innerLruCache,
		}},
	}
}

func newLruCache[K comparable, V any](c *Config) (localCache *LocalCache[K, V]) {
	// 初始化缓存实例 缓存容量默认256MB 最小512KB 最大8GB
	// The entry size need less than 1/1024 of cache size
	once.Do(func() {
		size, err := ParseSize(cfg.GetString("jupiter.cache.size"))
		if err != nil {
			xlog.Jupiter().Warn("localCache StdConfig ParseSize err", zap.Error(err))
		}
		if size < 512*KB || size > 8*GB {
			size = 256 * MB
		}
		innerCache = freecache.NewCache(int(size))
	})
	return &LocalCache[K, V]{
		cache: cache[K, V]{&myFreeCache[K, V]{
			config: c,
			cache:  innerCache,
		}},
	}
}
