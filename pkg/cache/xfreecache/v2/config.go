package xfreecache

import (
	"fmt"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Config struct {
	Cache         *freecache.Cache // 本地缓存实例【必填】
	Expire        time.Duration    // 失效时间 【必填】
	DisableMetric bool             // metric上报 false 开启  ture 关闭【选填，默认开启】
	Name          string           // 本地缓存名称，用于日志标识&metric上报【选填】
}

var (
	innerCache *freecache.Cache
	once       sync.Once
)

// DefaultConfig 返回默认配置 缓存容量256MB
func DefaultConfig() Config {
	once.Do(func() {
		innerCache = freecache.NewCache(int(256 * MB))
	})
	return Config{
		Cache:         innerCache,
		Expire:        2 * time.Minute,
		DisableMetric: false,
		Name:          "default",
	}
}

// SetCache 设置缓存实例
func (c Config) SetCache(cache *freecache.Cache) Config {
	c.Cache = cache
	return c
}

// SetExpire 设置缓存失效时间
func (c Config) SetExpire(t time.Duration) Config {
	c.Expire = t
	return c
}

// SetDisableMetric 设置是否进行metric上报
func (c Config) SetDisableMetric(disableMetric bool) Config {
	c.DisableMetric = disableMetric
	return c
}

// SetName 设置缓存名称
func (c Config) SetName(name string) Config {
	c.Name = name
	return c
}

// NewLocalCache 构建本地缓存实例 The entry size need less than 1/1024 of cache size
func NewLocalCache[K comparable, V any](c Config) (localCache *LocalCache[K, V]) {
	if c.Cache == nil {
		xlog.Jupiter().Panic("localCache NewLocalCache Cache nil", zap.Any("req", c))
	}
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache NewLocalCache expire err", zap.Any("req", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}

	localCache = &LocalCache[K, V]{
		cache: cache[K, V]{&localStorage{c}},
	}
	return
}
