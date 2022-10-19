package redisgo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/douyu/jupiter/pkg/xtrace"
	"github.com/douyu/jupiter/pkg/xtrace/jaeger"

	"github.com/douyu/jupiter/pkg/conf"

	"github.com/go-redis/redis/v8"
)

func Test_Interceptor(t *testing.T) {
	config := DefaultConfig()
	config.Addr = addr
	config.name = "test"
	t.Run("slow log", func(t *testing.T) {
		config.SlowLogThreshold = time.Nanosecond * 10
		client := config.BuildStub()
		client.Set(context.Background(), "redigo", "hello", time.Second)
		client.Pipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
			pipeliner.Del(context.Background(), "redigo")
			pipeliner.Get(context.Background(), "redigo")

			return nil
		})
		time.Sleep(time.Millisecond)
		client.Close()
	})

	t.Run("debug", func(t *testing.T) {
		config.Debug = true
		client := config.BuildStub()
		client.Set(context.Background(), "redigo", "hello", time.Second)
		client.Del(context.Background(), "redigo")
		client.Get(context.Background(), "redigo")

		time.Sleep(time.Millisecond)

		client.Pipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(context.Background(), "redigo", "hello", time.Second)
			pipeliner.Del(context.Background(), "redigo")
			pipeliner.Get(context.Background(), "redigo")
			return nil
		})
		client.Get(context.Background(), "redigo")

		client.Close()
	})

	t.Run("access", func(t *testing.T) {
		conf.Set("jupiter.trace.jaeger", map[string]interface{}{
			"addr": "wsd-jaeger-agent-go.pub.unp.oyw:6831",
			"rate": 1,
		})
		var con = jaeger.RawConfig("jupiter.trace.jaeger")
		xtrace.SetGlobalTracer(con.Build())
		ctx, span := xtrace.NewTracer(trace.SpanKindServer).Start(context.Background(), "test", nil)
		fmt.Println(span.SpanContext().TraceID())

		config.EnableAccessLogInterceptor = true
		client := config.BuildStub()

		client.Set(ctx, "redigo", "hello", time.Second)
		client.Del(ctx, "redigo")
		client.Get(ctx, "redigo")

		time.Sleep(time.Millisecond)

		client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(ctx, "redigo", "hello", time.Second)
			pipeliner.Del(ctx, "redigo")
			pipeliner.Get(ctx, "redigo")
			return nil
		})
		client.Get(ctx, "redigo")

		client.Close()
	})
	t.Run("trace", func(t *testing.T) {
		conf.Set("jupiter.trace.jaeger", map[string]interface{}{
			"addr": "wsd-jaeger-agent-go.pub.unp.oyw:6831",
			"rate": 1,
		})
		var con = jaeger.RawConfig("jupiter.trace.jaeger")
		xtrace.SetGlobalTracer(con.Build())

		config.EnableTraceInterceptor = true
		client := config.BuildStub()
		ctx := context.Background()
		client.Set(ctx, "redigo", "hello", time.Second)
		client.Del(ctx, "redigo")
		client.Get(ctx, "redigo")

		time.Sleep(time.Millisecond)

		client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(ctx, "redigo", "hello", time.Second)
			pipeliner.Del(ctx, "redigo")
			pipeliner.Get(ctx, "redigo")
			return nil
		})
		client.Get(ctx, "redigo")

		client.Close()
	})
}
