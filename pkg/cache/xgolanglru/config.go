package xgolanglru

import (
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"time"

	"github.com/douyu/jupiter/pkg/cache"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Config struct {
	Expire        time.Duration `json:"expire" toml:"expire"`               // 失效时间 【必填】
	DisableMetric bool          `json:"disableMetric" toml:"disableMetric"` // metric上报 false 开启  ture 关闭【选填，默认开启】
	Size          int           `json:"size" toml:"size"`                   // 缓存大小【选填，默认200000】
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
	key := "jupiter.xgolanglru." + name
	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		xlog.Jupiter().Warn("localCache StdConfig unmarshal config",
			zap.Error(err), zap.Any("key", key))
	}
	config.Name = name
	return config
}

// StdNew 构建本地缓存实例 The entry size need less than 1/1024 of cache size
func StdNew[K comparable, V any](name string) (localCache *cache.Cache[K, V]) {
	c := StdConfig(name)
	return New[K, V](c)
}

// New 构建本地缓存实例 缓存容量默认256MB 最小512KB 最大8GB
func New[K comparable, V any](c *Config) (localCache *cache.Cache[K, V]) {
	// 校验参数
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache New expire err", zap.Any("config", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}
	// 如果没有配置大小，默认200000
	if c.Size == 0 {
		c.Size = 200000
	}
	return &cache.Cache[K, V]{
		Storage: &localStorage[K, V]{
			config: c,
			Cache:  expirable.NewLRU[string, V](c.Size, nil, c.Expire),
		},
	}
}
