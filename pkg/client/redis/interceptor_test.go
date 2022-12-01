package redis

import (
	"context"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/core/xtrace/jaeger"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func Test_Interceptor(t *testing.T) {
	config := DefaultConfig()
	config.Master.Addr = addr
	config.Slaves.Addr = []string{addr2}
	config.name = "test"
	t.Run("slow log", func(t *testing.T) {
		config.SlowLogThreshold = time.Nanosecond * 10
		client, _ := config.Build()
		client.CmdOnMaster().Set(context.Background(), "redigo", "hello", time.Second)
		client.CmdOnSlave().Set(context.Background(), "redigo", "hello", time.Second)

		client.CmdOnMaster().Pipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
			pipeliner.Del(context.Background(), "redigo")
			pipeliner.Get(context.Background(), "redigo")

			return nil
		})
		time.Sleep(time.Millisecond)
		client.CmdOnMaster().Close()
	})

	t.Run("debug", func(t *testing.T) {
		config.Debug = true
		client, _ := config.Build()

		client.CmdOnMaster().Set(context.Background(), "redigo", "hello", time.Second)
		client.CmdOnMaster().Del(context.Background(), "redigo")
		client.CmdOnSlave().Get(context.Background(), "redigo")

		time.Sleep(time.Millisecond)

		client.CmdOnMaster().Pipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(context.Background(), "redigo", "hello", time.Second)
			pipeliner.Del(context.Background(), "redigo")
			pipeliner.Get(context.Background(), "redigo")
			return nil
		})
		client.CmdOnMaster().Get(context.Background(), "redigo")

	})

	t.Run("access", func(t *testing.T) {
		ctx := context.TODO()
		config.EnableAccessLogInterceptor = true
		client, _ := config.Build()

		client.CmdOnMaster().Set(ctx, "redigo", "hello", time.Second)
		client.CmdOnMaster().Del(ctx, "redigo")
		client.CmdOnMaster().Get(ctx, "redigo")

		time.Sleep(time.Millisecond)

		client.CmdOnMaster().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(ctx, "redigo", "hello", time.Second)
			pipeliner.Del(ctx, "redigo")
			pipeliner.Get(ctx, "redigo")
			return nil
		})
		client.CmdOnMaster().Get(ctx, "redigo")

	})
	t.Run("trace", func(t *testing.T) {
		xtrace.SetGlobalTracer((&jaeger.Config{
			Name:     "trace",
			Endpoint: "localhost:6831",
			Sampler:  1,
		}).Build())

		config.EnableTraceInterceptor = true
		client, _ := config.Build()
		ctx := context.Background()
		client.CmdOnMaster().Set(ctx, "redigo", "hello", time.Second)
		client.CmdOnMaster().Del(ctx, "redigo")
		client.CmdOnMaster().Get(ctx, "redigo")

		time.Sleep(time.Millisecond)

		client.CmdOnMaster().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(ctx, "redigo", "hello", time.Second)
			pipeliner.Del(ctx, "redigo")
			pipeliner.Get(ctx, "redigo")
			return nil
		})
		client.CmdOnMaster().Get(ctx, "redigo")

	})
	t.Run("sentinel", func(t *testing.T) {

		client, _ := config.Build()
		ctx := context.Background()
		assert.Equal(t, redis.Nil, client.CmdOnMaster().Get(ctx, "redigo").Err())

		time.Sleep(time.Millisecond)

		_, err := client.CmdOnMaster().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Get(ctx, "redigo")
			return nil
		})
		assert.Equal(t, redis.Nil, err)
		assert.Contains(t, client.CmdOnMaster().Do(ctx, "get").Err().Error(), "wrong number of arguments")

		_, err = client.CmdOnMaster().Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
			pipeliner.Do(ctx, "get")
			return nil
		})
		assert.Contains(t, err, "wrong number of arguments")
	})
}
