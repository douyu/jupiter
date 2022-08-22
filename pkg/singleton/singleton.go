package singleton

import (
	"sync"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/spf13/cast"
)

var singleton sync.Map

func Load(module constant.Module, key string) (interface{}, bool) {
	return singleton.Load(cast.ToString(module) + key)
}

func Store(module constant.Module, key string, val interface{}) {
	singleton.Store(cast.ToString(module)+key, val)
}
