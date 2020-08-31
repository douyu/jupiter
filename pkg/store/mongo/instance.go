package mongo

import (
	"sync"

	"github.com/globalsign/mgo"
)

var _instances = sync.Map{}

// Range 遍历所有实例
func Range(fn func(name string, db *mgo.Session) bool) {
	_instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*mgo.Session))
	})
}

// Get ...
func Get(name string) *mgo.Session {
	if ins, ok := _instances.Load(name); ok {
		return ins.(*mgo.Session)
	}
	return nil
}
