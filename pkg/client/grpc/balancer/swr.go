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

package balancer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/smallnest/weighted"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	// NameSmoothWeightRoundRobin ...
	NameSmoothWeightRoundRobin = "swr"
)

// PickerBuildInfo ...
type PickerBuildInfo struct {
	// ReadySCs is a map from all ready SubConns to the Addresses used to
	// create them.
	ReadySCs map[balancer.SubConn]base.SubConnInfo
	*attributes.Attributes
}

// PickerBuilder ...
type PickerBuilder interface {
	Build(info PickerBuildInfo) balancer.V2Picker
}

func init() {
	balancer.Register(
		NewBalancerBuilderV2(NameSmoothWeightRoundRobin, &swrPickerBuilder{}, base.Config{HealthCheck: true}),
	)
}

type swrPickerBuilder struct{}

// Build ...
func (s swrPickerBuilder) Build(info PickerBuildInfo) balancer.V2Picker {
	return newSWRPicker(info)
}

type swrPicker struct {
	readySCs     map[balancer.SubConn]base.SubConnInfo
	mu           sync.Mutex
	next         int
	buckets      *weighted.SW
	routeBuckets map[string]*weighted.SW
	*attributes.Attributes
}

func newSWRPicker(info PickerBuildInfo) *swrPicker {
	picker := &swrPicker{
		buckets:      &weighted.SW{},
		readySCs:     info.ReadySCs,
		routeBuckets: map[string]*weighted.SW{},
	}
	picker.parseBuildInfo(info)
	return picker
}

// Pick ...
func (p *swrPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var buckets = p.buckets
	if bs, ok := p.routeBuckets[info.FullMethodName]; ok {
		// 根据URI进行流量分组路由
		buckets = bs
	}

	sub, ok := buckets.Next().(balancer.SubConn)
	if ok {
		return balancer.PickResult{SubConn: sub}, nil
	}

	return balancer.PickResult{}, errors.New("pick failed")
}

func (p *swrPicker) parseBuildInfo(info PickerBuildInfo) {
	var hostedSubConns = map[string]balancer.SubConn{}
	var groupedSubConns = map[string][]balancer.SubConn{}

	for subConn, info := range info.ReadySCs {
		p.buckets.Add(subConn, 1)
		if info.Address.Attributes != nil {
			if serviceInfo, ok := info.Address.Attributes.Value(constant.KeyServiceInfo).(server.ServiceInfo); ok {
				// todo(gorexlv): 分组
				group := serviceInfo.Group
				if _, ok := groupedSubConns[group]; !ok {
					groupedSubConns[group] = make([]balancer.SubConn, 0)
				}
				groupedSubConns[group] = append(groupedSubConns[group], subConn)
			}
		}
		host := info.Address.Addr
		hostedSubConns[host] = subConn
		p.buckets.Add(subConn, 1)
	}

	if info.Attributes == nil {
		return
	}

	providerConfig, ok := info.Attributes.Value(constant.KeyProviderConfig).(registry.ProviderConfig)
	if !ok {
		return
	}

	fmt.Printf("providerConfig => %v\n", providerConfig)

	consumerConfigs, ok := info.Attributes.Value(constant.KeyConsumerConfig).(map[string]registry.ConsumerConfig)
	if !ok {
		return
	}

	fmt.Printf("consumerConfigs => %v\n", consumerConfigs)

	// 路由配置
	routeConfigs, ok := info.Attributes.Value(constant.KeyRouteConfig).(map[string]registry.RouteConfig)
	if !ok {
		return
	}
	for _, config := range routeConfigs {
		if _, ok := p.routeBuckets[config.URI]; !ok {
			p.routeBuckets[config.URI] = &weighted.SW{}
		}

		// 基于Group的权重配置, 同一分组下的IP分配同一个权重值
		for group, weight := range config.Upstream.Groups {
			sConns, ok := groupedSubConns[group]
			if !ok {
				continue
			}

			for _, sConn := range sConns {
				p.routeBuckets[config.URI].Add(sConn, weight)
			}
		}

		// 基于Node IP的权重配置, 如果配置了对应Node，将会覆盖Group中配置的权重
		for node, weight := range config.Upstream.Nodes {
			if sConn, ok := hostedSubConns[node]; ok {
				p.routeBuckets[config.URI].Add(sConn, weight)
			}
		}

	}
}
