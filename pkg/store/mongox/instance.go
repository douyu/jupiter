package mongox

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var _instances = sync.Map{}

// Range 遍历所有实例
func Range(fn func(name string, db *mongo.Client) bool) {
	_instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*mongo.Client))
	})
}

// Get ...
func Get(name string) *mongo.Client {
	if ins, ok := _instances.Load(name); ok {
		return ins.(*mongo.Client)
	}
	return nil
}
