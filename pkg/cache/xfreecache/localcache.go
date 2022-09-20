package xfreecache

import (
	"github.com/coocood/freecache"
	prome "github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type LocalCache struct {
	cache
}

type localStorage struct {
	cache *freecache.Cache
	req   Config
}

func (l *localStorage) SetCacheData(key string, data []byte) (err error) {
	err = l.cache.Set([]byte(key), data, int(l.req.Expire.Seconds()))
	// metric report
	if !l.req.DisableMetric {
		prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "EntryCount").Set(float64(l.cache.EntryCount()))
	}
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
	data, err = l.cache.Get([]byte(key))
	// metric report
	if !l.req.DisableMetric {
		prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "HitRate").Set(l.cache.HitRate())
		prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "HitCount").Set(float64(l.cache.HitCount()))
		prome.CacheHandleGauge.WithLabelValues(prome.TypeLocalCache, l.req.Name, "MissCount").Set(float64(l.cache.MissCount()))
	}
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
