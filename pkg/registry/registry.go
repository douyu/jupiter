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
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/douyu/jupiter/pkg/server"
)

// Event ...
type Event uint8

const (
	// EventUnknown ...
	EventUnknown Event = iota
	// EventUpdate ...
	EventUpdate
	// EventDelete ...
	EventDelete
)

// Kind ...
type Kind uint8

const (
	// KindUnknown ...
	KindUnknown Kind = iota
	// KindProvider ...
	KindProvider
	// KindConfigurator ...
	KindConfigurator
	// KindConsumer ...
	KindConsumer
)

// String ...
func (kind Kind) String() string {
	switch kind {
	case KindProvider:
		return "providers"
	case KindConfigurator:
		return "configurators"
	case KindConsumer:
		return "consumers"
	default:
		return "unknown"
	}
}

// ToKind ...
func ToKind(kindStr string) Kind {
	switch kindStr {
	case "providers":
		return KindProvider
	case "configurators":
		return KindConfigurator
	case "consumers":
		return KindConsumer
	default:
		return KindUnknown
	}
}

// ServerInstance ...
type ServerInstance struct {
	Scheme string
	IP     string
	Port   int
	Labels map[string]string
}

// EventMessage ...
type EventMessage struct {
	Event
	Kind
	Name    string
	Scheme  string
	Address string
	Message interface{}
}

// Registry register/unregister service
// registry impl should control rpc timeout
type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	UnregisterService(context.Context, *server.ServiceInfo) error
	ListServices(context.Context, string, string) ([]*server.ServiceInfo, error)
	WatchServices(context.Context, string, string) (chan Endpoints, error)
	io.Closer
}

//GetServiceKey ..
func GetServiceKey(prefix string, s *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/%s/%s://%s", prefix, s.Name, s.Kind.String(), s.Scheme, s.Address)
}

//GetServiceValue ..
func GetServiceValue(s *server.ServiceInfo) string {
	val, _ := json.Marshal(s)
	return string(val)
}

//GetService ..
func GetService(s string) *server.ServiceInfo {
	var si server.ServiceInfo
	json.Unmarshal([]byte(s), &si)
	return &si
}

// Nop registry, used for local development/debugging
type Nop struct{}

// ListServices ...
func (n Nop) ListServices(ctx context.Context, s string, s2 string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Nop) WatchServices(ctx context.Context, s string, s2 string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Nop) RegisterService(context.Context, *server.ServiceInfo) error { return nil }

// UnregisterService ...
func (n Nop) UnregisterService(context.Context, *server.ServiceInfo) error { return nil }

// Close ...
func (n Nop) Close() error { return nil }

// Configuration ...
type Configuration struct {
	Routes []Route           `json:"routes"` // 配置客户端路由策略
	Labels map[string]string `json:"labels"` // 配置服务端标签: 分组
}

// Route represents route configuration
type Route struct {
	// 路由方法名
	Method string `json:"method" toml:"method"`
	// 路由权重组, 按比率在各个权重组中分配流量
	WeightGroups []WeightGroup `json:"weightGroups" toml:"weightGroups"`
	// 路由部署组, 将流量导入部署组
	Deployment string `json:"deployment" toml:"deployment"`
}

// WeightGroup ...
type WeightGroup struct {
	Group  string `json:"group" toml:"group"`
	Weight int    `json:"weight" toml:"weight"`
}

// AddressList ...
// type AddressList struct {
// 	serverName string
//
// 	// regInfos map[string]*server.ServiceInfo
// 	// cfgInfos map[string]*server.ConfigInfo
//
// 	// TODO 2019/9/10 gorexlv: need lock
// 	regInfos map[string]*server.ServiceInfo // 注册信息
// 	cfgInfos map[string]*RouteConfig        // 配置信息
//
// 	mtx sync.RWMutex
// }
//
// func NewAddressList(serverName string) *AddressList {
// 	return &AddressList{
// 		serverName: serverName,
// 		regInfos:   make(map[string]*url.Values),
// 		cfgInfos:   make(map[string]*url.Values),
// 	}
// }
//
// func (al *AddressList) String() string {
// 	return ""
// }
//
// func (al *AddressList) List2() {
//
// }
//
// // List ...
// func (al *AddressList) List() []resolver.Address {
// 	// TODO 2019/9/10 gorexlv:
// 	addrs := make([]resolver.Address, 0)
// 	al.mtx.RLock()
// 	defer al.mtx.RUnlock()
// 	for addr, values := range al.regInfos {
// 		metadata := *values
// 		address := resolver.Address{
// 			Addr:       addr,
// 			ServerName: al.serverName,
// 			Attributes: attributes.New(),
// 		}
// 		// if infos, ok := al.cfgInfos[addr]; ok {
// 		// 	for cfg := range *infos {
// 		// 		metadata.Set(cfg, infos.Get(cfg))
// 		// 		// address.Attributes.WithValues(cfg, infos.Get(cfg))
// 		// 	}
// 		// }
// 		if enable := metadata.Get("enable"); enable != "" && enable != "true" {
// 			continue
// 		}
// 		if metadata.Get("weight") == "" {
// 			metadata.Set("weight", "100")
// 		}
// 		// group: 客户端配置的分组，默认为default
// 		// metadata.Get("group"): 服务端配置的分组
// 		// if metadata.Get("group") != "" && metadata.Get("group") != group {
// 		// 	continue
// 		// }
// 		address.Metadata = &metadata
// 		addrs = append(addrs, address)
// 	}
// 	return addrs
// }
