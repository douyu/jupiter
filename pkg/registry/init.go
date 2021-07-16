// Copyright 2021 rex lv
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

package registry

import (
	"log"
	"sync"

	"github.com/douyu/jupiter/pkg/conf"
)

var _registerers = sync.Map{}
var registryBuilder = make(map[string]Builder)

type Config map[string]struct {
	Kind          string `json:"kind" description:"底层注册器类型, eg: etcdv3, consul"`
	ConfigKey     string `json:"configKey" description:"底册注册器的配置键"`
	DeplaySeconds int    `json:"deplaySeconds" description:"延迟注册"`
}

func init() {
	// 初始化注册中心
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init registry")
		var config Config
		if err := c.UnmarshalKey("jupiter.registry", &config); err != nil {
			log.Printf("hook config, read registry config failed: %v", err)
			return
		}

		for name, item := range config {
			var itemKind = item.Kind
			if itemKind == "" {
				itemKind = "etcdv3"
			}
			build, ok := registryBuilder[itemKind]
			if !ok {
				log.Printf("invalid registry kind: %s", itemKind)
				continue
			}
			_registerers.Store(name, build(item.ConfigKey))
			log.Printf("build registrerer %s with config: %s", name, item.ConfigKey)
		}
	})
}

type Builder func(string) Registry

func RegisterBuilder(kind string, build Builder) {
	if _, ok := registryBuilder[kind]; ok {
		log.Panicf("duplicate register registry builder: %s", kind)
	}
	registryBuilder[kind] = build
}
