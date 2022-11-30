package singleton

import (
	"sync"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/spf13/cast"
)

var singleton sync.Map

func genkey(module constant.Module, key string) string {
	return cast.ToString(int(module)) + key
}

func Load(module constant.Module, key string) (interface{}, bool) {
	return singleton.Load(genkey(module, key))
}

func Store(module constant.Module, key string, val interface{}) {
	singleton.Store(genkey(module, key), val)
}
