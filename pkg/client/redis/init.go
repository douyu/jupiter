package redis

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"

	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/server/governor"
)

var instances = sync.Map{}

type storeRedis struct {
	ClientCluster *redis.ClusterClient
	ClientStub    *redis.Client
}

func init() {
	governor.HandleFunc("/debug/redis/stats", func(w http.ResponseWriter, r *http.Request) {
		_ = jsoniter.NewEncoder(w).Encode(stats())
	})
	go monitor()
}
func monitor() {
	for {
		instances.Range(func(key, val interface{}) bool {
			name := key.(string)
			obj := val.(*storeRedis)
			var poolStats *redis.PoolStats
			if obj.ClientStub != nil {
				poolStats = obj.ClientStub.PoolStats()
			}
			if obj.ClientCluster != nil {
				poolStats = obj.ClientCluster.PoolStats()
			}

			if poolStats != nil {
				prome.ClientStatsGauge.Set(float64(poolStats.Hits), prome.TypeRedis, name, "hits")
				prome.ClientStatsGauge.Set(float64(poolStats.Misses), prome.TypeRedis, name, "misses")
				prome.ClientStatsGauge.Set(float64(poolStats.Timeouts), prome.TypeRedis, name, "timeouts")
				prome.ClientStatsGauge.Set(float64(poolStats.TotalConns), prome.TypeRedis, name, "total_conns")
				prome.ClientStatsGauge.Set(float64(poolStats.IdleConns), prome.TypeRedis, name, "idle_conns")
				prome.ClientStatsGauge.Set(float64(poolStats.StaleConns), prome.TypeRedis, name, "stale_conns")
			}
			return true
		})
		time.Sleep(time.Second * 10)
	}
}

// stats
func stats() (stats map[string]interface{}) {
	stats = make(map[string]interface{})
	instances.Range(func(key, val interface{}) bool {
		name := key.(string)
		obj := val.(*storeRedis)
		if obj.ClientStub != nil {
			stats[name] = obj.ClientStub.PoolStats()
		}
		if obj.ClientCluster != nil {
			stats[name] = obj.ClientCluster.PoolStats()
		}
		return true
	})
	return
}
