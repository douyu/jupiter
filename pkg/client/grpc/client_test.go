package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/util/xtest/server/yell"
	"github.com/douyu/jupiter/proto/testproto"
	"github.com/stretchr/testify/assert"
)

// TestBase test direct dial with New()
func TestDirectGrpc(t *testing.T) {
	t.Run("test direct grpc", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		res, err := directClient.SayHello(ctx, &testproto.HelloRequest{
			Name: "hello",
		})
		assert.Nil(t, err)
		assert.Equal(t, res.Message, yell.RespFantasy.Message)
	})
}

func TestConfigBlockTrue(t *testing.T) {
	t.Run("test no address no block", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.DialTimeout = time.Second
		cfg.Debug = true
		conn := cfg.MustSingleton()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		res, err := testproto.NewGreeterClient(conn).SayHello(ctx, &testproto.HelloRequest{
			Name: "hello",
		})

		assert.ErrorContains(t, err, "missing address")
		assert.Nil(t, res)
	})
}

func TestAsyncConnect(t *testing.T) {
	t.Run("test async connect", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Addr = "127.0.0.1:9530"
		conn := cfg.Build()

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		res, err := testproto.NewGreeterClient(conn).SayHello(ctx, &testproto.HelloRequest{
			Name: "hello",
		})
		assert.NotNil(t, err)
		assert.Nil(t, res)

		go func() {
			startServer("127.0.0.1:9530", "test-async-server")
		}()

		assert.Eventually(t, func() bool {

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			fmt.Println(conn.GetState())
			res, err := testproto.NewGreeterClient(conn).SayHello(ctx, &testproto.HelloRequest{
				Name: "hello",
			})
			fmt.Println(err, res)
			return err == nil && res != nil
		}, 5*time.Second, time.Second)

	})
}
