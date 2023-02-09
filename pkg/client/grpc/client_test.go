package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestBase test direct dial with New()
func TestDirectGrpc(t *testing.T) {
	t.Run("test direct grpc", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		res, err := directClient.SayHello(ctx, &helloworldv1.SayHelloRequest{
			Name: "hello",
		})
		assert.Nil(t, err)
		assert.Equal(t, res.Data.Name, "hello")
	})
}

func TestConfigBlockTrue(t *testing.T) {
	t.Run("test no address no block", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.DialTimeout = time.Second
		cfg.Debug = true
		conn, err := cfg.Build()

		assert.NotNil(t, err)
		assert.Nil(t, conn)
		assert.Equal(t, "failed to build resolver: passthrough: received empty target in Build()", err.Error())
	})
}

func TestAsyncConnect(t *testing.T) {
	t.Run("test async connect", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Addr = "127.0.0.1:9530"
		conn := lo.Must(cfg.Build())

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		res, err := helloworldv1.NewGreeterServiceClient(conn).SayHello(ctx, &helloworldv1.SayHelloRequest{
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
			res, err := helloworldv1.NewGreeterServiceClient(conn).SayHello(ctx, &helloworldv1.SayHelloRequest{
				Name: "hello",
			})
			fmt.Println(err, res)
			return err == nil && res != nil
		}, 5*time.Second, time.Second)

	})
}
