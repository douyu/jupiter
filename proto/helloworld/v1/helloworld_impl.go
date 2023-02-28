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

func (s *FooServer) SayGoodBye(ctx context.Context, in *SayGoodByeRequest) (out *SayGoodByeResponse, err error) {
	var fm = new(SayGoodByeRequest_FieldMask)
	if in.Type == Type_TYPE_FILTER {
		fm = in.FieldMaskFilter()
	} else {
		fm = in.FieldMaskPrune()
	}
	out = &SayGoodByeResponse{
		Error: 0,
		Msg:   "请求正常",
		Data: &SayGoodByeResponse_Data{
			Age:  1,
			Name: "",
			Other: &OtherHelloMessage{
				Id:      1,
				Address: "bar",
			},
		},
	}
	if fm.MaskedInName() {
		out.Data.Name = in.GetName()
	}

	if fm.MaskedInAge() {
		out.Data.Age = in.GetAge()
	}
	_ = fm.Mask(out)
	return
}
