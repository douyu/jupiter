package xfreecache

import (
	"fmt"
	"time"

	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Config struct {
	CacheType     string        `json:"cacheType" toml:"cacheType"`         // 本地缓存类型 例如：freecache、lru等，默认使用freeCache
	Expire        time.Duration `json:"expire" toml:"expire"`               // 失效时间 【必填】
	DisableMetric bool          `json:"disableMetric" toml:"disableMetric"` // metric上报 false 开启  ture 关闭【选填，默认开启】
	Name          string        `json:"-" toml:"-"`                         // 本地缓存名称，用于日志标识&metric上报【选填】
}

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

type LocalCache[K comparable, V any] struct {
	cache[K, V]
}

// StdNew 构建本地缓存实例
func StdNew[K comparable, V any](name string) (localCache *LocalCache[K, V]) {
	c := StdConfig(name)
	return New[K, V](c)
}

// New 构建本地缓存实例
func New[K comparable, V any](c *Config) (localCache *LocalCache[K, V]) {
	// 校验参数
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache New expire err", zap.Any("config", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}

	// 根据配置选择不同的缓存类型
	switch c.CacheType {
	case cacheTypeFreeCache:
		localCache = newFreeCache[K, V](c)
	case cacheTypeLruCache:
		localCache = newLruCache[K, V](c)
	default:
		localCache = newFreeCache[K, V](c)
	}

	return
}
