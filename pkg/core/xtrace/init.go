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

package xtrace

import (
	"log"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/xtrace/jaeger"
	"github.com/douyu/jupiter/pkg/core/xtrace/otelgrpc"
)

func init() {
	// 加载完配置，初始化trace
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init trace config")

		key := constant.ConfigKey("trace.jaeger")
		if conf.Get(key) != nil {
			var config = jaeger.RawConfig(key)
			SetGlobalTracer(config.Build())

			return
		}

		key = constant.ConfigKey("trace.otelgrpc")
		if conf.Get(key) != nil {
			var config = otelgrpc.RawConfig(key)
			SetGlobalTracer(config.Build())

			return
		}
	})
}
