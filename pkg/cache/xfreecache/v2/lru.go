package xfreecache

import (
	"fmt"
	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	glru "github.com/hnlq715/golang-lru"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"reflect"
)

type myLru[K comparable, V any] struct {
	config *Config
	cache  *glru.ARCCache // lru本地缓存实例【必填】
}

func (l *myLru[K, V]) getCacheMap(key string, ids []K) (v map[K]V, idsNone []K) {
	var zero V
	v = make(map[K]V)
	idsNone = make([]K, 0, len(ids))

	// id去重
	ids = lo.Uniq(ids)
	for _, id := range ids {
		cacheKey := l.getKey(key, id)
		resT, ok := l.cache.Get(cacheKey)
		if ok {
			if value, ok1 := resT.(V); ok1 {
				if !reflect.DeepEqual(value, zero) {
					v[id] = value
				}
			}
		} else {
			idsNone = append(idsNone, id)
		}
		// metric report
		if !l.config.DisableMetric {
			if ok {
				prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "HitCount", "").Inc()
			} else {
				prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "MissCount", "").Inc()
			}
		}
	}
	return
}

func (l *myLru[K, V]) setCacheMap(key string, idsNone []K, fn func([]K) (map[K]V, error), v map[K]V) (err error) {
	args := []zap.Field{zap.Any("key", key), zap.Any("ids", idsNone)}

	if len(idsNone) == 0 {
		return
	}

	// 执行函数
	resMap, err := fn(idsNone)
	if err != nil {
		xlog.Jupiter().Error("setCacheMap doMap", append(args, zap.Error(err))...)
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
		var cacheData V
		if val, ok := resMap[id]; ok {
			cacheData = val
		}

		cacheKey := l.getKey(key, id)
		l.cache.AddEx(cacheKey, cacheData, l.config.Expire)
		// metric report
		if !l.config.DisableMetric {
			prome.CacheHandleCounter.WithLabelValues(prome.TypeLocalCache, l.config.Name, "SetSuccess", "").Inc()
		}
	}
	return
}

func (l *myLru[K, V]) getKey(key string, id K) string {
	return fmt.Sprintf("%s:%v", key, id)
}
