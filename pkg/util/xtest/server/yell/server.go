package yell

import (
	"context"
	"errors"
	"time"

	"github.com/douyu/jupiter/proto/testproto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FooServer ...
type FooServer struct {
	name string
	hook func(context.Context)
}

// SetName ...
func (s *FooServer) SetName(f string) {
	s.name = f
}

// SetHook ...
func (s *FooServer) SetHook(f func(context.Context)) {
	s.hook = f
}

// ErrFoo ...
var ErrFoo = errors.New("error foo")

// RespFantasy ...
var RespFantasy = &testproto.SayHelloResponse{Data: &testproto.SayHelloResponse_Data{Name: "fantasy"}}

// RespBye ...
var RespBye = &testproto.SayHelloResponse{Data: &testproto.SayHelloResponse_Data{Name: "test"}}

// StatusFoo ...
var StatusFoo = status.Errorf(codes.DataLoss, ErrFoo.Error())

// SayHello ...
func (s *FooServer) SayHello(ctx context.Context, in *testproto.SayHelloRequest) (out *testproto.SayHelloResponse, err error) {
	// sleep to test cost time
	time.Sleep(20 * time.Millisecond)
	switch in.Name {
	case "traceHook":
		s.hook(ctx)
		err = StatusFoo
	case "needErr":
		err = StatusFoo
	case "slow":
		time.Sleep(500 * time.Millisecond)
		out = RespFantasy
	case "needPanic":
		panic("go dead!")
	default:
		out = RespFantasy
	}
	return
}

// // StreamHello ...
// func (s *FooServer) StreamHello(ss testproto.Greeter_StreamHelloServer) (err error) {

// 	for {
// 		in, _ := ss.Recv()
// 		switch in.Name {
// 		case "bye":
// 			return ss.Send(RespBye)
// 		case "needErr":
// 			return StatusFoo
// 		default:
// 			return ss.Send(RespFantasy)
// 		}
// 	}
// }

// // StreamHello ...
// func (s *FooServer) WhoServer(ctx context.Context, in *testproto.WhoServerReq) (out *testproto.WhoServerReply, err error) {
// 	return &testproto.WhoServerReply{Message: s.name}, nil
// }
