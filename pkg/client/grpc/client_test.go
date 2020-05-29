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

package grpc

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/resolver"

	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/client/grpc/balancer"
	etcdv3_registry "github.com/douyu/jupiter/pkg/registry/etcdv3"
)

/*
/jupiter/main/providers/grpc://127.0.0.1:20102
*/

func init() {
	resolver.Register(etcdv3_registry.Config{
		ReadTimeout: time.Second * 3,
		Prefix:      "jupiter",
		Config: &etcdv3.Config{
			Endpoints:      []string{"127.0.0.1:2379"},
			ConnectTimeout: time.Second * 3,
			Secure:         false,
		},
	}.BuildResolver())
}

func Test_newClient(t *testing.T) {
	config := DefaultConfig()
	config = config.WithDialOption(grpc.WithInsecure(), grpc.WithBalancerName(balancer.NameSmoothWeightRoundRobin))
	config.Address = "etcd:///main"
	client := helloworld.NewGreeterClient(config.Build())
	ctx, ctxErr := context.WithTimeout(context.Background(), time.Second*3)
	defer ctxErr()
	rep, err := client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "hi",
	})

	t.Logf("errj => %+v\n", err)
	t.Logf("rep => %+v\n", rep)
	time.Sleep(time.Second)
}
