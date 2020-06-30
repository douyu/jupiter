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

	"github.com/douyu/jupiter/pkg/util/xrand"
	"github.com/smallnest/weighted"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	// NameSmoothWeightRoundRobin ...
	NameSmoothWeightRoundRobin = "swr"
)

func newGWRBuilder(policy string) balancer.Builder {
	return base.NewBalancerBuilderV2(policy, &groupPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newGWRBuilder(NameSmoothWeightRoundRobin))
}

type groupPickerBuilder struct {
	policy string
}

// Build ...
func (gpb *groupPickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}

	return newWeightPicker(info.ReadySCs)
}

type weightGroup struct {
	name    string
	buckets *weighted.RandW
}

type weightPicker struct {
	mu      sync.Mutex
	next    int
	buckets *weighted.SW
}

func newWeightPicker(readySCs map[balancer.SubConn]base.SubConnInfo) *weightPicker {
	wp := &weightPicker{
		next:    xrand.Intn(len(readySCs)),
		buckets: &weighted.SW{},
	}
	/*
		if group name is provided by client, check all ready sub connection:
		1. if these is no attribute, drop  it
		2. if no group info in attribute, drop it
		3. if no weight info in attribute, drop it
		if no group name provided, degrade to round_robin balance
	*/
	var groups map[string]*attributes.Attributes
	for subConn, info := range readySCs {
		attributes := info.Address.Attributes
		if attributes == nil {
			continue
		}

		config, ok := attributes.Value("meta").(*Config)
		if !ok {
			continue
		}

		if config.Group == "" || !config.Enable {
			continue
		}

	}

	for group, attributes := range groups {
		fmt.Printf("group = %+v\n", group)
		wp.buckets.Add(attributes, weight)
	}

	return wp
}

// Pick ...
func (p *weightPicker) Pick(opts balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	sub, ok := p.bucket.Next().(balancer.SubConn)
	if ok {
		return balancer.PickResult{
			SubConn: sub,
		}, nil
	}

	return balancer.PickResult{}, errors.New("pick failed")
}

// Config ...
type Config struct {
	Group  string
	Weight int
	Enable bool
}
