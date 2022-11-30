package impl

import (
	"context"

	"github.com/douyu/jupiter/proto/testproto"
)

type TestProjectImp struct {
	testproto.UnimplementedGreeterServer
}

func (s *TestProjectImp) SayHello(ctx context.Context, req *testproto.HelloRequest) (*testproto.HelloReply, error) {
	return &testproto.HelloReply{
		Message: "hello",
	}, nil
}
