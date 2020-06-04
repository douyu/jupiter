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
