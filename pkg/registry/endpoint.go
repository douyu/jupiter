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

package registry

import (
	"encoding/json"

	"github.com/douyu/jupiter/pkg/server"
)

// Endpoints ...
type Endpoints struct {
	// 服务节点列表
	Nodes map[string]server.ServiceInfo

	// 路由配置
	RouteConfigs map[string]RouteConfig

	// 消费者元数据
	ConsumerConfigs map[string]ConsumerConfig

	// 服务元信息
	ProviderConfigs map[string]ProviderConfig
}

func newEndpoints() *Endpoints {
	return &Endpoints{
		Nodes:           make(map[string]server.ServiceInfo),
		RouteConfigs:    make(map[string]RouteConfig),
		ConsumerConfigs: make(map[string]ConsumerConfig),
		ProviderConfigs: make(map[string]ProviderConfig),
	}
}

func (in *Endpoints) DeepCopy() *Endpoints {
	if in == nil {
		return nil
	}

	out := newEndpoints()
	in.DeepCopyInfo(out)
	return out
}

func (in *Endpoints) DeepCopyInfo(out *Endpoints) {
	for key, info := range in.Nodes {
		out.Nodes[key] = info
	}
	for key, config := range in.RouteConfigs {
		out.RouteConfigs[key] = config
	}
	for key, config := range in.ConsumerConfigs {
		out.ConsumerConfigs[key] = config
	}
	for key, config := range in.ProviderConfigs {
		out.ProviderConfigs[key] = config
	}
}

// ProviderConfig config of provider
// 通过这个配置，修改provider的属性
type ProviderConfig struct {
	ID     string `json:"id"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`

	Region     string            `json:"region"`
	Zone       string            `json:"zone"`
	Deployment string            `json:"deployment"`
	Metadata   map[string]string `json:"metadata"`
	Enable     bool              `json:"enable"`
}

// ConsumerConfig config of consumer
// 客户端调用app的配置
type ConsumerConfig struct {
	ID     string `json:"id"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

// RouteConfig ...
type RouteConfig struct {
	ID     string `json:"id" toml:"id"`
	Scheme string `json:"scheme" toml:"scheme"`
	Host   string `json:"host" toml:"host"`

	Deployment string   `json:"deployment"`
	URI        string   `json:"uri"`
	Upstream   Upstream `json:"upstream"`
}

// String ...
func (config RouteConfig) String() string {
	bs, _ := json.Marshal(config)
	return string(bs)
}

// Upstream represents upstream balancing config
type Upstream struct {
	Nodes  map[string]int `json:"nodes"`
	Groups map[string]int `json:"groups"`
}
