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
	"sync"
	"testing"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

var weightNodes = map[string]int{
	"127.0.0.1:9091": 1,
	"127.0.0.1:9092": 2,
	"127.0.0.1:9093": 3,
	"127.0.0.1:9094": 4,
	"127.0.0.1:9095": 5,
	"127.0.0.1:9096": 6,
	"127.0.0.1:9097": 7,
	"127.0.0.1:9098": 8,
	"127.0.0.1:9099": 9,
}

var groupedNodes = map[string]int {
	"red": 1,
	"green": 2,
}

func Test_swrPicker(t *testing.T) {
	buildInfo:=PickerBuildInfo{
		ReadySCs:   map[balancer.SubConn]base.SubConnInfo{},
	}

	t.Run("with node upstream attributes", func(t *testing.T){
		buildInfo.Attributes = attributes.New(
			constant.KeyRouteConfig,
			map[string]registry.RouteConfig{
				"/routes/1": {
					URI:        "/hello",
					Upstream:   registry.Upstream{ Nodes:  weightNodes },
				},
			},
		)

		t.Run("pick first", func(t *testing.T) {
			address := "127.0.0.1:9091"
			subConn := &mockSubConn{addr:address}
			buildInfo.ReadySCs[subConn] = base.SubConnInfo{
				Address: resolver.Address{
					Addr:       address,
					ServerName: "testing",
				},
			}
			result, err := newSWRPicker(buildInfo).Pick(balancer.PickInfo{
				FullMethodName: "/hello",
				Ctx:            context.Background(),
			})
			assert.Nil(t, err)
			assert.NotNil(t, result)
		})

		t. Run("weight pick", func(t *testing.T) {
			buildInfo.ReadySCs = map[balancer.SubConn]base.SubConnInfo{}
			for addr := range weightNodes {
				subConn := &mockSubConn{addr:addr}
				buildInfo.ReadySCs[subConn] = base.SubConnInfo{
					Address: resolver.Address{
						Addr:       addr,
						ServerName: "testing",
					},
				}
			}
			var picker = newSWRPicker(buildInfo)
			t.Run("grouped route", func(t *testing.T) {
				var nodeCount = map[string]int{}
				for i:=0;i<45;i++{
					result, err := picker.Pick(balancer.PickInfo{
						FullMethodName: "/hello",
						Ctx:            context.Background(),
					})
					assert.Nil(t, err)
					assert.NotNil(t, result)
					nodeCount[result.SubConn.(*mockSubConn).addr]++
				}
				assert.Equal(t, nodeCount, weightNodes)
			})
			t.Run("ungrouped route", func(t *testing.T) {
				for i:=0;i<45;i++{
					result, err := picker.Pick(balancer.PickInfo{
						FullMethodName: "/ungrouped_route",
						Ctx:            context.Background(),
					})
					assert.Nil(t, err)
					assert.NotNil(t, result)
				}
			})
		})
	})

	t.Run("no attributes", func(t *testing.T){
		buildInfo.Attributes = nil
		t.Run("pick first", func(t *testing.T) {
			address := "127.0.0.1:9091"
			subConn := &mockSubConn{addr:address}
			buildInfo.ReadySCs[subConn] = base.SubConnInfo{
				Address: resolver.Address{
					Addr:       address,
					ServerName: "testing",
				},
			}
			result, err := newSWRPicker(buildInfo).Pick(balancer.PickInfo{
				FullMethodName: "/hello",
				Ctx:            context.Background(),
			})
			assert.Nil(t, err)
			assert.NotNil(t, result)
		})
	})

}

type mockSubConn struct {
	addr string
	balancer.SubConn
}

type mockClientConn struct {
	balancer.ClientConn

	mu       sync.Mutex
	subConns map[balancer.SubConn]resolver.Address
}

func newMockClientConn() *mockClientConn {
	return &mockClientConn{
		subConns: make(map[balancer.SubConn]resolver.Address),
	}
}

func (mcc *mockClientConn) NewSubConn(addrs []resolver.Address, opts balancer.NewSubConnOptions) (balancer.SubConn, error) {
	sc := &mockSubConn{}
	mcc.mu.Lock()
	defer mcc.mu.Unlock()
	mcc.subConns[sc] = addrs[0]
	return sc, nil
}

func (mcc *mockClientConn) RemoveSubConn(sc balancer.SubConn) {
	mcc.mu.Lock()
	defer mcc.mu.Unlock()
	delete(mcc.subConns, sc)
}
