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

	"github.com/douyu/jupiter/pkg/trace"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/xlog"
)

// run: go run main.go -config=config.toml
type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func main() {
	app := NewEngine()

	for k := 0; k < 10; k++ {
		time.Sleep(time.Second)
		traceTest()

	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func traceTest() {
	// 1. 从配置文件中初始化
	process1 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process1")
		defer span.Finish()
		// todo something
		time.Sleep(time.Second)
		fmt.Println("finish", "process1")

	}

	process2 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process2")
		defer span.Finish()
		process1(ctx)
		time.Sleep(time.Second)
		fmt.Println("finish", "process2")

	}

	process3 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process3")
		defer span.Finish()
		process2(ctx)
		time.Sleep(time.Second)
		fmt.Println("finish", "process3")
	}

	process3(context.Background())
	return
}
