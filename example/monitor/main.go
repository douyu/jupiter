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
	"github.com/douyu/jupiter/example/grpc/helloworld/helloworld"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/xlog"
)

func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Error(err.Error())
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveGRPC,
		func() error {
			client := etcdv3.StdConfig("myetcd").MustBuild()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			// 添加数据
			_, err := client.Put(ctx, fmt.Sprintf("/prometheus/job/%s/%s/%s", "jupiter", "monitor-demo", "127.0.0.1:9999"), "127.0.0.1:9999")
			if err != nil {
				xlog.Panic(err.Error())
			}
			return nil
		},
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) serveGRPC() error {
	server := xgrpc.StdConfig("grpc").MustBuild()
	helloworld.RegisterGreeterServer(server.Server, new(Greeter))
	return eng.Serve(server)
}

type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: "Hello Jupiter",
	}, nil
}
