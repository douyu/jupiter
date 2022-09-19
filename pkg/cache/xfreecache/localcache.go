package xfreecache

import (
	"fmt"
	"time"

	"github.com/coocood/freecache"
	prome "github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type LocalCache struct {
	Cache
}

type LocalCacheReq struct {
	Size          Size          // 缓存容量,最小512*1024 【必填】
	Expire        time.Duration // 失效时间 【必填】
	DisableMetric bool          // metric上报 【选填，默认开启】
	Name          string        // 本地缓存名称，用于日志标识&metric上报【选填】
}

// NewLocalCache 新建本地缓存实例 The entry size need less than 1/1024 of cache size
func NewLocalCache(req LocalCacheReq) (cache *LocalCache) {
	if req.Size < 512*KB {
		xlog.Jupiter().Panic("cache NewLocalCache size err", zap.Any("req", req))
	}
	if req.Expire == 0 {
		xlog.Jupiter().Panic("cache NewLocalCache expire err", zap.Any("req", req))
	}
	if len(req.Name) == 0 {
		req.Name = fmt.Sprintf("cache-%d", time.Now().UnixNano())
	}

	cacheLocal := freecache.NewCache(int(req.Size))
	cache = &LocalCache{
		Cache: Cache{&localStorage{cacheLocal, req}},
	}
	return
}

type localStorage struct {
	cache *freecache.Cache
	req   LocalCacheReq
}

func (l *localStorage) SetCacheData(key string, data []byte) (err error) {
	// metric上报
	if !l.req.DisableMetric {
		defer func() {
			prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "EntryCount").Set(float64(l.cache.EntryCount()))
		}()
	}

	err = l.cache.Set([]byte(key), data, int(l.req.Expire.Seconds()))
	if err != nil {
		xlog.Jupiter().Error("cache SetCacheData", zap.String("data", string(data)), zap.Error(err))
		if err == freecache.ErrLargeEntry || err == freecache.ErrLargeKey {
			err = nil
		}
		return
	}
	return
}

func (l *localStorage) GetCacheData(key string) (data []byte, err error) {
	// metric上报
	if !l.req.DisableMetric {
		defer func() {
			prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "MissCount").Set(float64(l.cache.MissCount()))
			prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "HitRate").Set(l.cache.HitRate())
			prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "AverageAccessTime").Set(float64(l.cache.AverageAccessTime() * 1000))
		}()
	}

	data, err = l.cache.Get([]byte(key))
	if err == freecache.ErrNotFound || data == nil {
		err = nil
		return
	}
	if err != nil {
		xlog.Jupiter().Error("cache GetCacheData", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
