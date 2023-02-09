package helloworldv1

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestGreeterServiceGRPC_SayHello_0(t *testing.T) {
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
			name: "case 1: gRPC with protobuf",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *grpc.Server {
					echo := grpc.NewServer()

					return echo
				},
			},
			args: args{
				createReq: func() *http.Request {
					body, _ := proto.Marshal(&SayHelloRequest{Name: "bob"})
					hdr, body := msgHeader(body, nil)
					data := bytes.NewBuffer(hdr)
					data.Write(body)
					req := httptest.NewRequest(
						"POST", "http://localhost/helloworld.v1.GreeterService/SayHello",
						data,
					)
					req.ProtoMajor = 2
					req.Header.Add("Content-Type", "application/grpc")
					return req
				},
			},
			wantErr: false,
			wantRes: &SayHelloResponse{
				Data: &SayHelloResponse_Data{
					Name: "bob",
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

// msgHeader returns a 5-byte header for the message being transmitted and the
// payload, which is compData if non-nil or data otherwise.
func msgHeader(data, compData []byte) (hdr []byte, payload []byte) {
	hdr = make([]byte, 5)
	if compData != nil {
		hdr[0] = byte(1)
		data = compData
	} else {
		hdr[0] = byte(0)
	}

	// Write length of payload into buf
	binary.BigEndian.PutUint32(hdr[1:], uint32(len(data)))
	return hdr, data
}
