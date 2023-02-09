package xfreecache

import (
	"fmt"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Config struct {
	Cache         *freecache.Cache `json:"-" toml:"-"`                         // 本地缓存实例【必填】
	Expire        time.Duration    `json:"expire" toml:"expire"`               // 失效时间 【必填】
	DisableMetric bool             `json:"disableMetric" toml:"disableMetric"` // metric上报 false 开启  ture 关闭【选填，默认开启】
	Name          string           `json:"-" toml:"-"`                         // 本地缓存名称，用于日志标识&metric上报【选填】
}

var (
	innerCache *freecache.Cache
	once       sync.Once
)

// DefaultConfig 返回默认配置 缓存容量256MB 最小512KB 最大8GB
func DefaultConfig(size Size) *Config {
	once.Do(func() {
		if size < 512*KB || size > 8*GB {
			size = 256 * MB
		}
		innerCache = freecache.NewCache(int(size))
	})
	return &Config{
		Cache:         innerCache,
		Expire:        2 * time.Minute,
		DisableMetric: false,
		Name:          "default",
	}
}

// StdConfig 返回标准配置
func StdConfig(name string) *Config {
	size, err := ParseSize(cfg.GetString("jupiter.cache.size"))
	if err != nil {
		xlog.Jupiter().Error("localCache StdConfig ParseSize err", zap.Error(err))
	}

	config := DefaultConfig(size)
	key := "jupiter.cache." + name
	if err = cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		xlog.Jupiter().Warn("localCache StdConfig unmarshal config",
			zap.Error(err), zap.Any("key", key))
	}
	config.Name = name
	return config
}

// StdNew 构建本地缓存实例 The entry size need less than 1/1024 of cache size
func StdNew[K comparable, V any](name string) (localCache *LocalCache[K, V]) {
	c := StdConfig(name)
	return New[K, V](c)
}

func New[K comparable, V any](c *Config) (localCache *LocalCache[K, V]) {
	if c.Cache == nil {
		xlog.Jupiter().Panic("localCache New Cache nil", zap.Any("config", c))
	}
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache New expire err", zap.Any("config", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}

	localCache = &LocalCache[K, V]{
		cache: cache[K, V]{&localStorage{c}},
	}
	return
}
