package xfreecache

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type storage interface {
	// SetCacheData 设置缓存数据 key：缓存key data：缓存数据
	SetCacheData(key string, data []byte) (err error)
	// GetCacheData 存储缓存数据 key：缓存key data：缓存数据
	GetCacheData(key string) (data []byte, err error)
}

type cache[K comparable, V any] struct {
	storage
}

// GetAndSetCacheData 获取缓存后数据
func (c *cache[K, V]) GetAndSetCacheData(key string, id K, fn func() (V, error)) (value V, err error) {
	resMap, err := c.GetAndSetCacheMap(key, []K{id}, func([]K) (map[K]V, error) {
		innerVal, innerErr := fn()
		return map[K]V{id: innerVal}, innerErr
	})
	value = resMap[id]
	return
}

// GetAndSetCacheMap 获取缓存后数据 map形式
func (c *cache[K, V]) GetAndSetCacheMap(key string, ids []K, fn func([]K) (map[K]V, error)) (v map[K]V, err error) {
	args := []zap.Field{zap.Any("key", key), zap.Any("ids", ids)}

	v = make(map[K]V)

	// id去重
	ids = lo.Uniq(ids)
	idsNone := make([]K, 0, len(ids))
	pool := getPool[V]()
	for _, id := range ids {
		cacheKey := c.getKey(key, id)
		if resT, innerErr := c.GetCacheData(cacheKey); innerErr == nil && resT != nil {
			var value V
			// 反序列化
			value, err = unmarshalWithPool[V](resT, pool)
			if err != nil {
				return
			}
			v[id] = value
			continue
		}
		idsNone = append(idsNone, id)
	}

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
	for k, value := range resMap {
		v[k] = value
	}

	// 写入缓存
	for _, id := range idsNone {
		var (
			cacheData V
			data      []byte
		)

		if val, ok := v[id]; ok {
			cacheData = val
		}
		// 序列化
		data, err = marshal(cacheData)

		if err != nil {
			xlog.Jupiter().Error("GetAndSetCacheMap Marshal", append(args, zap.Error(err))...)
			return
		}

		cacheKey := c.getKey(key, id)
		err = c.SetCacheData(cacheKey, data)
		if err != nil {
			xlog.Jupiter().Error("GetAndSetCacheMap setCacheData", append(args, zap.Error(err))...)
			return
		}
	}
	return
}

func (c *cache[K, V]) getKey(key string, id K) string {
	return fmt.Sprintf("%s:%v", key, id)
}
