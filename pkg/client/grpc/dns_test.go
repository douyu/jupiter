package grpc

import (
	"context"
	"testing"

	"github.com/douyu/jupiter/proto/testproto/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDNS(t *testing.T) {
	config := DefaultConfig()
	config.Addr = "dns:///localhost:9528"
	config.Debug = true
	cc := testproto.NewGreeterServiceClient(lo.Must(config.Build()))

	res, err := cc.SayHello(context.Background(), &testproto.SayHelloRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
