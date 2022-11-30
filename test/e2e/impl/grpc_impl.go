package impl

import (
	"context"

	"github.com/douyu/jupiter/proto/testproto/v1"
)

type TestProjectImp struct {
	testproto.UnimplementedGreeterServiceServer
}

func (s *TestProjectImp) SayHello(ctx context.Context, req *testproto.SayHelloRequest) (*testproto.SayHelloResponse, error) {
	return &testproto.SayHelloResponse{
		Data: &testproto.SayHelloResponse_Data{
			Name: req.Name,
		},
	}, nil
}
