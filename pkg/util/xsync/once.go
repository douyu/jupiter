package xsync

import "sync"

var onceMapper = sync.Map{}

func Once(name string) *sync.Once {
	var value, ok = onceMapper.Load(name)
	if ok {
		return value.(*sync.Once)
	}
	newValue := &sync.Once{}
	onceMapper.Store(name, newValue)
	return newValue
}
