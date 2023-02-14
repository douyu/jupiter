package benchmark

import (
	"bytes"
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douyu/jupiter/pkg/core/encoding"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/util/xerror"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func BenchmarkServer(b *testing.B) {

	b.Run("gRPC with protobuf", func(b *testing.B) {
		server := grpc.NewServer()
		impl := new(impl)
		helloworldv1.RegisterGreeterServiceServer(server, impl)
		body, _ := proto.Marshal(&helloworldv1.SayHelloRequest{Name: "bob"})
		hdr, body := msgHeader(body, nil)
		data := bytes.NewBuffer(hdr)
		data.Write(body)

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(
				"POST", "http://localhost/helloworld.v1.GreeterService/SayHello",
				data,
			)
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with reflect gRPC", func(b *testing.B) {
		server := echo.New()
		impl := new(impl)
		server.POST("/v1/helloworld.Greeter/SayHello", xecho.GRPCProxyWrapper(impl.SayHello))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with echo json", func(b *testing.B) {
		server := echo.New()
		server.POST("/v1/helloworld.Greeter/SayHello", echoJsonHandler)

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with grpc gateway", func(b *testing.B) {

		mux := runtime.NewServeMux()
		helloworldv1.RegisterGreeterServiceHandlerServer(context.TODO(), mux, new(impl))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with echo stdjson gateway", func(b *testing.B) {
		server := echo.New()

		helloworldv1.RegisterGreeterServiceEchoServer(server, new(impl))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with echo protojson gateway", func(b *testing.B) {
		server := echo.New()
		server.JSONSerializer = new(encoding.ProtoJsonSerializer)

		helloworldv1.RegisterGreeterServiceEchoServer(server, new(impl))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

	b.Run("HTTP with gin stdjson gateway", func(b *testing.B) {
		gin.SetMode(gin.ReleaseMode)
		server := gin.New()

		helloworldv1.RegisterGreeterServiceGinServer(server, new(impl))

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/v1/helloworld.Greeter/SayHello", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			// fmt.Println(rec)
			// b.Fail()
		}
	})

}

func echoJsonHandler(c echo.Context) error {
	req := new(helloworldv1.SayHelloRequest)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	if req.Name != "bob" {
		return c.JSON(http.StatusOK, &helloworldv1.SayHelloResponse{
			Error: uint32(xerror.InvalidArgument.GetEcode()),
			Msg:   "invalid name",
			Data:  &helloworldv1.SayHelloResponse_Data{},
		})
	}
	return c.JSON(http.StatusOK, &helloworldv1.SayHelloResponse{
		Msg: "",
		Data: &helloworldv1.SayHelloResponse_Data{
			Name: "hello bob",
		},
	})
}

type impl struct {
	helloworldv1.UnimplementedGreeterServiceServer
}

func (*impl) SayHello(ctx context.Context, req *helloworldv1.SayHelloRequest) (*helloworldv1.SayHelloResponse, error) {
	if req.Name != "bob" {
		return &helloworldv1.SayHelloResponse{
			Error: uint32(xerror.InvalidArgument.GetEcode()),
			Msg:   "invalid name",
			Data:  &helloworldv1.SayHelloResponse_Data{},
		}, nil
	}
	return &helloworldv1.SayHelloResponse{
		Msg: "",
		Data: &helloworldv1.SayHelloResponse_Data{
			Name: "hello bob",
		},
	}, nil
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
