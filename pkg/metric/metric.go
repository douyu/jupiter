// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"time"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	prometheus.Registerer
	prometheus.Gatherer

	BulkRegister(...prometheus.Collector) error
}

var (
	// TypeHTTP ...
	TypeHTTP = "http"
	// TypeGRPCUnary ...
	TypeGRPCUnary = "unary"
	// TypeGRPCStream ...
	TypeGRPCStream = "stream"
	// TypeRedis ...
	TypeRedis = "redis"
	TypeGorm  = "gorm"
	// TypeRocketMQ ...
	TypeRocketMQ = "rocketmq"
	// TypeWebsocket ...
	TypeWebsocket = "ws"

	// TypeMySQL ...
	TypeMySQL = "mysql"

	// CodeJob
	CodeJobSuccess = "ok"
	// CodeJobFail ...
	CodeJobFail = "fail"
	// CodeJobReentry ...
	CodeJobReentry = "reentry"

	// CodeCache
	CodeCacheMiss = "miss"
	// CodeCacheHit ...
	CodeCacheHit = "hit"
)

var (
	// ServerHandleCounter ...
	ServerHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "server_handle_total",
		Labels:    []string{"type", "method", "client", "code"},
	}.Build()

	// ServerHandleHistogram ...
	ServerHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "server_handle_seconds",
		Labels:    []string{"type", "method", "client"},
	}.Build()

	// ClientHandleCounter ...
	ClientHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "client_handle_total",
		Labels:    []string{"type", "name", "method", "server", "code"},
	}.Build()

	// ClientHandleHistogram ...
	ClientHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "client_handle_seconds",
		Labels:    []string{"type", "name", "method", "server"},
	}.Build()

	// JobHandleCounter ...
	JobHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "job_handle_total",
		Labels:    []string{"type", "name", "code"},
	}.Build()

	// JobHandleHistogram ...
	JobHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "job_handle_seconds",
		Labels:    []string{"type", "name"},
	}.Build()

	LibHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_seconds",
		Labels:    []string{"type", "method", "address"},
	}.Build()
	// LibHandleCounter ...
	LibHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_total",
		Labels:    []string{"type", "method", "address", "code"},
	}.Build()

	LibHandleSummary = SummaryVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "lib_handle_stats",
		Labels:    []string{"name", "status"},
	}.Build()

	// CacheHandleCounter ...
	CacheHandleCounter = CounterVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "cache_handle_total",
		Labels:    []string{"type", "name", "method", "code"},
	}.Build()

	// CacheHandleHistogram ...
	CacheHandleHistogram = HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "cache_handle_seconds",
		Labels:    []string{"type", "name", "method"},
	}.Build()

	// BuildInfoGauge ...
	BuildInfoGauge = GaugeVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      "build_info",
		Labels:    []string{"name", "id", "env", "region", "zone", "version", "go_version"},
		// Labels:    []string{"name", "aid", "mode", "region", "zone", "app_version", "jupiter_version", "start_time", "build_time", "go_version"},
	}.Build()

	// LogLevelCounter ...
	LogLevelCounter = NewCounterVec("log_level_total", []string{"name", "lv"})
)

func init() {
	conf.OnLoaded(func(c *conf.Configuration) {
		BuildInfoGauge.WithLabelValues(
			pkg.Name(),
			pkg.AppID(),
			pkg.AppMode(),
			pkg.AppRegion(),
			pkg.AppZone(),
			pkg.AppVersion(),
			pkg.GoVersion(),
		).Set(float64(time.Now().UnixNano() / 1e6))
	})
}
