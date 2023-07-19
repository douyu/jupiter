package xfreecache

import (
	"fmt"
	"github.com/coocood/freecache"
	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"reflect"
)

type myFreeCache[K comparable, V any] struct {
	config *Config
	cache  *freecache.Cache // 本地缓存实例【必填】
}

func (l *myFreeCache[K, V]) getCacheMap(key string, ids []K) (v map[K]V, idsNone []K) {
	var zero V
	v = make(map[K]V)
	idsNone = make([]K, 0, len(ids))

	// id去重
	ids = lo.Uniq(ids)
	for _, id := range ids {
		cacheKey := l.getKey(key, id)
		resT, innerErr := l.getCacheData(cacheKey)
		if innerErr == nil && resT != nil {
			var value V
			// 反序列化
			value, innerErr = unmarshal[V](resT)
			if innerErr != nil {
				xlog.Jupiter().Error("cache unmarshalWithPool", zap.String("key", key), zap.Error(innerErr))
			} else {
				if !reflect.DeepEqual(value, zero) {
					v[id] = value
				}
			}
		}
		if innerErr != nil {
			idsNone = append(idsNone, id)
		}
	}
	return
}

func (l *myFreeCache[K, V]) setCacheMap(key string, idsNone []K, fn func([]K) (map[K]V, error), v map[K]V) (err error) {
	args := []zap.Field{zap.Any("key", key), zap.Any("ids", idsNone)}

	if len(idsNone) == 0 {
		return
	}

	// 执行函数
	resMap, err := fn(idsNone)
	if err != nil {
		xlog.Jupiter().Error("GetAndSetCacheMap doMap", append(args, zap.Error(err))...)
		return
	}

	// 填入返回中
	if v != nil {
		for k, value := range resMap {
			v[k] = value
		}
	}

	// 写入缓存
	for _, id := range idsNone {
		var (
			cacheData V
			data      []byte
		)

		if val, ok := resMap[id]; ok {
			cacheData = val
		}
		// 序列化
		data, err = marshal(cacheData)

		if err != nil {
			xlog.Jupiter().Error("GetAndSetCacheMap Marshal", append(args, zap.Error(err))...)
			return
		}

		cacheKey := l.getKey(key, id)
		err = l.setCacheData(cacheKey, data)
		if err != nil {
			xlog.Jupiter().Error("GetAndSetCacheMap setCacheData", append(args, zap.Error(err))...)
			return
		}
	}
	return
}

func (l *myFreeCache[K, V]) getKey(key string, id K) string {
	return fmt.Sprintf("%s:%v", key, id)
}

func (l *myFreeCache[K, V]) setCacheData(key string, data []byte) (err error) {
	err = l.cache.Set([]byte(key), data, int(l.config.Expire.Seconds()))
	// metric report
	if !l.config.DisableMetric {
		if err != nil {
			prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "SetFail", err.Error()).Inc()
		} else {
			prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "SetSuccess", "").Inc()
		}
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

func (l *myFreeCache[K, V]) getCacheData(key string) (data []byte, err error) {
	data, err = l.cache.Get([]byte(key))
	// metric report
	if !l.config.DisableMetric {
		if err != nil {
			prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "MissCount", err.Error()).Inc()
		} else {
			prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "HitCount", "").Inc()
		}
	}
	if err != nil && err != freecache.ErrNotFound {
		xlog.Jupiter().Error("cache GetCacheData", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
