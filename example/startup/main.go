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
	"time"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/worker"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/labstack/echo/v4"
)

func main() {
	var app jupiter.Application
	app.Startup()
	app.Serve(startHTTPServer())
	app.Serve(startGRPCServer())
	app.Schedule(startWorker())
	app.Run()
}

func startHTTPServer() server.Server {
	server := xecho.DefaultConfig().Build()
	server.GET("/hello", func(ctx echo.Context) error {
		return ctx.JSON(200, "Gopher Wuhan")
	})
	return server
}

func startGRPCServer() server.Server {
	server := xgrpc.DefaultConfig().Build()
	// helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
	return server
}

func startWorker() worker.Worker {
	cron := xcron.DefaultConfig().Build()
	cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(func() error {
		return nil
	}))
	return cron
}
