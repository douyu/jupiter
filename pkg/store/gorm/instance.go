package gorm

import (
	"sync"
)

var instances = sync.Map{}

// Range 遍历所有实例
func Range(fn func(name string, db *DB) bool) {
	instances.Range(func(key, val interface{}) bool {
		return fn(key.(string), val.(*DB))
	})
}

// Configs
func Configs() map[string]interface{} {
	var rets = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		return true
	})

	return rets
}

// Stats
func Stats() (stats map[string]interface{}) {
	stats = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		db := val.(*DB)

		stats[name] = db.DB().Stats()
		return true
	})

	return
}
