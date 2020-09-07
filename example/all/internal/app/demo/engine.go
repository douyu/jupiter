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

package demo

import (
	"time"

	sentinel_echo "github.com/alibaba/sentinel-golang/adapter/echo"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/example/all/internal/app/greeter"
	"github.com/douyu/jupiter/pkg/sentinel"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/examples/helloworld/helloworld"

	sentinel_grpc "github.com/alibaba/sentinel-golang/adapter/grpc"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}

	if err := eng.Startup(
		eng.startJobs,
		eng.initSentinel,
		xgo.ParallelWithError(
			eng.serveHTTP,
			eng.serveGRPC,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}

	return eng
}

func (eng *Engine) initSentinel() error {
	var config = sentinel.DefaultConfig()
	config.FlowRules = append(config.FlowRules, &flow.FlowRule{
		Resource:   "GET:/ping",
		MetricType: flow.QPS,
		Count:      1,
	})
	return config.Build()
}

func (eng *Engine) startJobs() error {
	cron := xcron.StdConfig("demo").Build()
	cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(eng.execJob))
	return eng.Schedule(cron)
}

func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").Build()
	server.Use(
		sentinel_echo.SentinelMiddleware(
			// customize resource extractor if required
			// method_path by default
			sentinel_echo.WithResourceExtractor(func(ctx echo.Context) string {
				return ctx.Request().Method + ":" + ctx.Path()
			}),
			// customize block fallback if required
			// abort with status 429 by default
			sentinel_echo.WithBlockFallback(func(ctx echo.Context) error {
				return ctx.JSON(400, map[string]interface{}{
					"err":  "too many requests; the quota used up",
					"code": 10222,
				})
			}),
		),
	)
	server.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(200, "pong")
	})
	server.GET("/panic", func(ctx echo.Context) error {
		panic("panic")
	})
	// this is a demo: support proxy for http to grpc controller
	g := greeter.Greeter{}
	server.GET("/grpc", xecho.GRPCProxyWrapper(g.SayHello))
	server.POST("/grpc-post", xecho.GRPCProxyWrapper(g.SayHello))
	return eng.Serve(server)
}

func (eng *Engine) serveGRPC() error {
	server := xgrpc.StdConfig("grpc").
		WithUnaryInterceptor(sentinel_grpc.NewUnaryServerInterceptor()).
		WithStreamInterceptor(sentinel_grpc.NewStreamServerInterceptor()).
		Build()

	helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
	return eng.Serve(server)
}

func (eng *Engine) execJob() error {
	xlog.Info("exec job", xlog.String("info", "print info"))
	xlog.Warn("exec job", xlog.String("warn", "print warning"))
	return nil
}
