package xfreecache

import (
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type storage interface {
	// SetCacheData 设置缓存数据 key：缓存key data：缓存数据
	SetCacheData(key string, data []byte) (err error)
	// GetCacheData 存储缓存数据 key：缓存key data：缓存数据
	GetCacheData(key string) (data []byte, err error)
}

type cache struct {
	storage
}

// GetAndSetCacheData 获取缓存后数据
func (c *cache) GetAndSetCacheData(key string, fn func() ([]byte, error)) (v []byte, err error) {
	v, err = c.GetCacheData(key)
	if err == nil && v != nil {
		return
	}
	// 执行程序
	v, err = fn()
	// 如果程序报错了就不进行缓存
	if err != nil {
		xlog.Jupiter().Error("cache SetAndGetCacheData do", zap.String("key", key), zap.Error(err))
		return
	}
	err = c.SetCacheData(key, v)
	if err != nil {
		return
	}
	return
}
