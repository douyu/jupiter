package grpc

import (
	"context"
	"testing"

	"github.com/douyu/jupiter/proto/testproto"
	"github.com/stretchr/testify/assert"
)

func TestDNS(t *testing.T) {
	config := DefaultConfig()
	config.Addr = "dns:///localhost:9528"
	config.Debug = true
	cc := testproto.NewGreeterClient(config.Build())

	res, err := cc.SayHello(context.Background(), &testproto.HelloRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
