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
	"net/http"

	"github.com/douyu/jupiter/pkg/govern"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TypeServerHttp ...
	TypeServerHttp = "http"
	// TypeServerUnary ...
	TypeServerUnary = "unary"
	// TypeServerStream ...
	TypeServerStream = "stream"
	// TypeLibRocketMq ...
	TypeLibRocketMq = "rocketMq"
)

var (
	// ServerHandleHistogram 指标: 服务类型，调用方法，客户端标识，返回的状态码
	ServerMetricsHandler = NewServerMetrics()

	// ClientHandleHistogram 指标: 客户端类型，客户端名称，调用方法，目标，返回的状态码
	ClientMetricsHandler = NewClientMetrics()

	// WorkerHandleHistogram 指标: 类型，任务名，执行状态码
	WorkerMetricsHandler = NewWorkerMetrics()
)

func init() {
	govern.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}
