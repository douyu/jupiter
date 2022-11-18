// Copyright 2022 Douyu
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

package sentinel

import (
	"github.com/alibaba/sentinel-golang/api"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/core/metric"
)

var (
	sentinelReqeust = metric.NewCounterVec("sentinel_request",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelSuccess = metric.NewCounterVec("sentinel_success",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelExceptionsThrown = metric.NewCounterVec("sentinel_exceptions_thrown",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelBlocked = metric.NewCounterVec("sentinel_blocked",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelRt = metric.NewHistogramVec("sentinel_rt",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelState = metric.NewGaugeVec("sentinel_state",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})

	WithError        = base.WithError
	WithResourceType = api.WithResourceType
	WithTrafficType  = api.WithTrafficType
)

const (
	language = "go"

	SENTINEL_DATASOURCE_ETCD  = "etcd"
	SENTINEL_DATASOURCE_FILES = "files"
)

type (
	SentinelEntry = base.SentinelEntry
	BlockError    = base.BlockError
	EntryContext  = base.EntryContext
	EntryOption   = sentinel.EntryOption
)

var (
	stdConfig Config
)

func init() {
	hooks.Register(hooks.Stage_AfterLoadConfig, func() {
		_ = build()
	})
}

// build 基于标准配置构建sentinel.
func build() error {
	stdConfig = StdConfig()

	return stdConfig.Build()
}

// Entry 执行熔断策略.
func Entry(resource string, opts ...EntryOption) (*SentinelEntry, *BlockError) {
	return stdConfig.Entry(resource, opts...)
}
