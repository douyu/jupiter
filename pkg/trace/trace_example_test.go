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

package trace_test

import (
	"context"
	"fmt"
	"time"

	"github.com/douyu/jupiter/pkg/trace"
)

func ExampleTraceFunc() {
	// 1. 从配置文件中初始化
	process1 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process1")
		defer span.Finish()

		// todo something
		fmt.Println("err", ctx.Err())
		time.Sleep(time.Second)
	}

	process2 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process2")
		defer span.Finish()
		process1(ctx)
	}

	process3 := func(ctx context.Context) {
		span, ctx := trace.StartSpanFromContext(ctx, "process3")
		defer span.Finish()
		process2(ctx)
	}

	process3(context.Background())
}
