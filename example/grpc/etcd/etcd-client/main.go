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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/grpc"
	"github.com/douyu/jupiter/pkg/client/grpc/balancer"
	"github.com/douyu/jupiter/pkg/client/grpc/resolver"
	"github.com/douyu/jupiter/pkg/registry/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Error(err.Error())
	}
	fmt.Printf("111 = %+v\n", 111)
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.initResolver,
		eng.consumer,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) initResolver() error {
	resolver.Register("etcd", etcdv3.StdConfig("wh").Build())
	return nil
}

func (eng *Engine) consumer() error {
	config := grpc.StdConfig("etcdserver")
	config.BalancerName = balancer.NameSmoothWeightRoundRobin

	client := helloworld.NewGreeterClient(config.Build())
	go func() {
		for {
			resp, err := client.SayHello(context.Background(), &helloworld.HelloRequest{
				Name: "jupiter",
			})
			if err != nil {
				fmt.Printf("err = %+v\n", err)
				xlog.Error(err.Error())
			} else {
				fmt.Printf("resp.Message = %+v\n", resp.Message)
				xlog.Info("receive response", xlog.String("resp", resp.Message))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
