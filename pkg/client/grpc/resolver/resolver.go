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

package resolver

import (
	"context"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Register ...
func Register(name string, reg registry.Registry) {
	resolver.Register(&baseBuilder{
		name: name,
		reg:  reg,
	})
}

type baseBuilder struct {
	name string
	reg  registry.Registry
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	endpoints, err := b.reg.WatchServices(context.Background(), target.Endpoint, "grpc")
	if err != nil {
		return nil, err
	}

	var stop = make(chan struct{})
	xgo.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.New(
						constant.KeyRouteConfig, endpoint.RouteConfigs, // 路由配置
						constant.KeyProviderConfig, endpoint.ProviderConfigs, // 服务提供方元信息
						constant.KeyConsumerConfig, endpoint.ConsumerConfigs, // 服务消费方配置信息
					),
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = target.Endpoint
					address.Attributes = attributes.New(constant.KeyServiceInfo, node)
					state.Addresses = append(state.Addresses, address)
				}
				cc.UpdateState(state)
			case <-stop:
				return
			}
		}
	})

	return &baseResolver{
		stop: stop,
	}, nil
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	stop chan struct{}
}

// ResolveNow ...
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {}

// Close ...
func (b *baseResolver) Close() { b.stop <- struct{}{} }
