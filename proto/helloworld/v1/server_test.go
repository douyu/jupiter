package helloworldv1

import (
	context "context"
	"errors"
	"time"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// FooServer ...
type FooServer struct {
	UnimplementedGreeterServiceServer

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

// StatusFoo ...
var StatusFoo = status.Errorf(codes.DataLoss, ErrFoo.Error())

// SayHello ...
func (s *FooServer) SayHello(ctx context.Context, in *SayHelloRequest) (out *SayHelloResponse, err error) {
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
		out = &SayHelloResponse{Data: &SayHelloResponse_Data{Name: in.Name}}
	case "needPanic":
		panic("go dead!")
	default:
		out = &SayHelloResponse{Data: &SayHelloResponse_Data{Name: in.Name}}
	}
	return
}
