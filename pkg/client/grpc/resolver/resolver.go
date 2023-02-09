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
	"strings"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/registry/etcdv3"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// NewEtcdBuilder returns a new etcdv3 resolver builder.
func NewEtcdBuilder(name string, registryConfig string) resolver.Builder {
	return &baseBuilder{
		name:           name,
		registryConfig: registryConfig,
	}
}

type baseBuilder struct {
	name string

	registryConfig string
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	reg := etcdv3.RawConfig(b.registryConfig).MustSingleton()

	serviceName := target.Endpoint()
	if !strings.HasSuffix(serviceName, "/") {
		serviceName += "/"
	}

	endpoints, err := reg.WatchServices(context.Background(), serviceName)
	if err != nil {
		xlog.Jupiter().Error("watch services failed", xlog.FieldErr(err))
		return nil, err
	}

	var stop = make(chan struct{})
	xgo.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				xlog.Jupiter().Debug("watch services finished", xlog.FieldValueAny(endpoint))

				var state = resolver.State{
					Addresses: make([]resolver.Address, 0),
					Attributes: attributes.
						New(constant.KeyRouteConfig, endpoint.RouteConfigs).             // 路由配置
						WithValue(constant.KeyProviderConfig, endpoint.ProviderConfigs). // 服务提供方元信息
						WithValue(constant.KeyConsumerConfig, endpoint.ConsumerConfigs), // 服务消费方配置信息,
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = serviceName
					address.Attributes = attributes.New(constant.KeyServiceInfo, node)
					state.Addresses = append(state.Addresses, address)
				}
				_ = cc.UpdateState(state)
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
