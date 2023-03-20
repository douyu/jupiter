package helloworldv1

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FooServer ...
type FooServer struct {
	UnimplementedGreeterServiceServer

	hook func(context.Context)
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
	var fm = new(SayHelloRequest_FieldMask)

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
	switch in.GetType() {
	case Type_TYPE_UNSPECIFIED:
		return
	case Type_TYPE_PRUNE:
		fm = in.FieldMaskPrune()
	case Type_TYPE_FILTER:
		fm = in.FieldMaskFilter()
	default:
		return
	}

	out = &SayHelloResponse{
		Error: 0,
		Msg:   "请求正常",
		Data: &SayHelloResponse_Data{
			Name:      "",
			AgeNumber: 18,
			Sex:       Sex_SEX_MALE,
			Metadata:  map[string]string{"Bar": "bar"},
		},
	}
	if fm.MaskedInName() {
		out.Data.Name = in.GetName()
	}
	_ = fm.Mask(out)
	return
}
