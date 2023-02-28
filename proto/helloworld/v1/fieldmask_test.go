package helloworldv1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestGreeterServiceGRPC_SayGoodBye_0(t *testing.T) {
	type fields struct {
		server       GreeterServiceServer
		createRouter func() *grpc.Server
	}
	type args struct {
		createReq func() *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantRes    proto.Message
		wantHeader http.Header
	}{
		{
			name: "case 1: fieldMask filter",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *grpc.Server {
					echo := grpc.NewServer()
					return echo
				},
			},
			args: args{
				createReq: func() *http.Request {
					protoreq := &SayGoodByeRequest{
						Name: "foo",
						Type: Type_TYPE_Filter,
					}
					protoreq.MaskInName().MaskOutDataName().MaskOutDataOther()

					fmt.Printf("request: %+v\n", protoreq.String())
					body, _ := proto.Marshal(protoreq)
					hdr, body := msgHeader(body, nil)
					data := bytes.NewBuffer(hdr)
					data.Write(body)
					req := httptest.NewRequest(
						"POST", "http://localhost/helloworld.v1.GreeterService/SayGoodBye",
						data,
					)
					req.ProtoMajor = 2
					req.Header.Add("Content-Type", "application/grpc")
					return req
				},
			},
			wantErr: false,
			wantRes: &SayGoodByeResponse{
				Error: 0,
				Msg:   "请求正常",
				Data: &SayGoodByeResponse_Data{
					Name: "foo",
					Other: &OtherHelloMessage{
						Id:      1,
						Address: "bar",
					},
				},
			},
		},
		{
			name: "case 2: fieldMask prune",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *grpc.Server {
					echo := grpc.NewServer()
					return echo
				},
			},
			args: args{
				createReq: func() *http.Request {
					protoreq := &SayGoodByeRequest{
						Name: "foo",
						Type: Type_TYPE_Prune,
					}
					protoreq.MaskInName().MaskOutDataName().MaskOutDataOther()

					fmt.Printf("request: %+v\n", protoreq.String())
					body, _ := proto.Marshal(protoreq)
					hdr, body := msgHeader(body, nil)
					data := bytes.NewBuffer(hdr)
					data.Write(body)
					req := httptest.NewRequest(
						"POST", "http://localhost/helloworld.v1.GreeterService/SayGoodBye",
						data,
					)
					req.ProtoMajor = 2
					req.Header.Add("Content-Type", "application/grpc")
					return req
				},
			},
			wantErr: false,
			wantRes: &SayGoodByeResponse{
				Error: 0,
				Msg:   "请求正常",
				Data: &SayGoodByeResponse_Data{
					Age: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := tt.fields.createRouter()
			RegisterGreeterServiceServer(router, tt.fields.server)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, tt.args.createReq())
			fmt.Printf("response: %+v\n", res.Body.String())
			hdr, body := msgHeader(lo.Must(proto.Marshal(tt.wantRes)), nil)
			data := bytes.NewBuffer(hdr)
			data.Write(body)

			assert.Equal(t, data.Bytes(), res.Body.Bytes())
			if len(tt.wantHeader) > 0 {
				assert.Equal(t, tt.wantHeader, res.Header())
			}
		})
	}
}
