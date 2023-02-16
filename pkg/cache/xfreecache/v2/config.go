package xfreecache

import (
	"fmt"
	"sync"
	"time"

	"github.com/coocood/freecache"
	cfg "github.com/douyu/jupiter/pkg/conf"
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

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Expire:        2 * time.Minute,
		DisableMetric: false,
		Name:          "default",
	}
}

// StdConfig 返回标准配置
func StdConfig(name string) *Config {
	config := DefaultConfig()
	key := "jupiter.cache." + name
	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
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

// New 构建本地缓存实例 缓存容量默认256MB 最小512KB 最大8GB
func New[K comparable, V any](c *Config) (localCache *LocalCache[K, V]) {
	// 校验参数
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache New expire err", zap.Any("config", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}
	if c.Cache == nil {
		// 初始化缓存实例
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
		c.Cache = innerCache
	}

	localCache = &LocalCache[K, V]{
		cache: cache[K, V]{&localStorage{c}},
	}
	return
}
