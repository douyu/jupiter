package singleton

import (
	"strconv"
	"sync"

	"github.com/douyu/jupiter/pkg/constant"
)

var singleton sync.Map

func Load(module constant.Module, key string) (interface{}, bool) {
	return singleton.Load(strconv.Itoa(int(module)) + key)
}

func Store(module constant.Module, key string, val interface{}) {
	singleton.Store(strconv.Itoa(int(module))+key, val)
}
