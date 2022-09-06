package cache

import (
	"github.com/douyu/jupiter/pkg/cache/xerr"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type Storage interface {
	// SetCacheData 设置缓存数据 key：缓存key data：缓存数据
	SetCacheData(key string, data []byte) (err error)
	// GetCacheData 存储缓存数据 key：缓存key data：缓存数据
	GetCacheData(key string) (data []byte, err error)
}

type Cache struct {
	Storage
}

// GetAndSetCacheData 获取缓存后数据
func (c *Cache) GetAndSetCacheData(key string, fn func() ([]byte, error)) (v []byte, err error) {
	v, err = c.GetCacheData(key)
	if err == nil && v != nil {
		return
	}
	// 执行程序
	v, err = c.do(fn)
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

// 执行程序
func (c *Cache) do(fn func() ([]byte, error)) (v []byte, err error) {
	normalReturn := false
	recovered := false

	// 使用双defer来区分来自runtime.Goexit的panic,
	defer func() {
		if !normalReturn && !recovered {
			err = xerr.ErrGoexit
		}

		if e, ok := err.(*xerr.PanicError); ok {
			panic(e)
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				if r := recover(); r != nil {
					err = xerr.NewPanicError(r)
				}
			}
		}()

		v, err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
	return
}
