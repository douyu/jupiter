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
	"context"
	"errors"
	"net/url"
	"strconv"
	"sync"

	"github.com/douyu/jupiter/pkg/util/xrand"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/smallnest/weighted"
)

const (
	// NameSmoothWeightRoundRobin ...
	NameSmoothWeightRoundRobin = "swr"
)

func newGWRBuilder(policy string) balancer.Builder {
	return base.NewBalancerBuilderWithConfig(policy, &groupPickerBuilder{
		policy: policy,
	}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newGWRBuilder(NameSmoothWeightRoundRobin))
}

type groupPickerBuilder struct {
	policy string
}

// Build ...
func (gpb *groupPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for _, sc := range readySCs {
		scs = append(scs, sc)
	}

	switch gpb.policy {
	case NameSmoothWeightRoundRobin:
		return newWeightPicker(readySCs)
	default:
		panic("balance invalid pick policy " + gpb.policy)
	}
}

type weightPicker struct {
	subConns []balancer.SubConn
	readySCs map[resolver.Address]balancer.SubConn

	mu     sync.Mutex
	next   int
	logger *xlog.Logger
	bucket *weighted.SW
}

func newWeightPicker(readySCs map[resolver.Address]balancer.SubConn) *weightPicker {
	wp := &weightPicker{
		readySCs: readySCs,
		next:     xrand.Intn(len(readySCs)),
		bucket:   &weighted.SW{},
	}

	for addr, sub := range readySCs {
		meta, ok := addr.Metadata.(*url.Values)
		if !ok {
			xlog.Error("metadata assert", xlog.Any("metadata", addr.Metadata))
			continue
		}
		// v1 版grpc没有weight字段，默认100
		weightStr := meta.Get("weight")
		if weightStr == "" {
			weightStr = "100"
		}

		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			xlog.Error("metadata weight", xlog.Any("metadata", addr.Metadata))
			continue
		}

		wp.bucket.Add(sub, weight)
	}

	return wp
}

// Pick ...
func (p *weightPicker) Pick(ctx context.Context, opts balancer.PickInfo) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	sub, ok := p.bucket.Next().(balancer.SubConn)
	if ok {
		return sub, nil, nil
	}

	return nil, nil, errors.New("pick failed")
}
