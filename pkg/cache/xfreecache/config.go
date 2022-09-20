package xfreecache

import (
	"fmt"
	"time"

	"github.com/coocood/freecache"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Config struct {
	Size          Size          // 缓存容量,最小512*1024 【必填】
	Expire        time.Duration // 失效时间 【必填】
	DisableMetric bool          // metric上报 false 开启  ture 关闭【选填，默认开启】
	Name          string        // 本地缓存名称，用于日志标识&metric上报【选填】
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Size:          256 * MB,
		Expire:        2 * time.Minute,
		DisableMetric: false,
		Name:          fmt.Sprintf("cache-%d", time.Now().UnixNano()),
	}
}

// Build 构建本地缓存实例 The entry size need less than 1/1024 of cache size
func (c Config) Build() (localCache *LocalCache) {
	if c.Size < 512*KB {
		xlog.Jupiter().Panic("localCache NewLocalCache size err", zap.Any("req", c))
	}
	if c.Expire == 0 {
		xlog.Jupiter().Panic("localCache NewLocalCache expire err", zap.Any("req", c))
	}
	if len(c.Name) == 0 {
		c.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}

	cacheLocal := freecache.NewCache(int(c.Size))
	localCache = &LocalCache{
		cache: cache{&localStorage{cacheLocal, c}},
	}
	return
}
