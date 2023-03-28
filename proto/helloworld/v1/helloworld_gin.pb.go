// Code generated by github.com/douyu/jupiter/cmd/protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v0.11.3
// - protoc             (unknown)

package helloworldv1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	metadata "google.golang.org/grpc/metadata"
	http "net/http"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the github.com/douyu/jupiter/cmd/protoc-gen-go-gin package it is being compiled against.
var _ = http.StatusOK
var _ = new(context.Context)
var _ = metadata.New
var _ = gin.Engine{}

type GreeterServiceGinServer interface {

	// Sends a hello greeting
	SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error)

	// Sends a hi greeting
	SayHi(context.Context, *SayHiRequest) (*SayHiResponse, error)
}

func RegisterGreeterServiceGinServer(r *gin.Engine, srv GreeterServiceGinServer) {
	s := _gin_GreeterService{
		server: srv,
		router: r,
	}
	s.registerService()
}

type _gin_GreeterService struct {
	server GreeterServiceGinServer
	router *gin.Engine
}

// Sends a hello greeting
func (s *_gin_GreeterService) _gin_handler_SayHello_0(ctx *gin.Context) {
	var in SayHelloRequest
	if err := ctx.ShouldBindUri(&in); err != nil {
		ctx.Error(err)
		return
	}
	if err := ctx.ShouldBind(&in); err != nil {
		ctx.Error(err)
		return
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request.Context(), md)
	out, err := s.server.(GreeterServiceGinServer).SayHello(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, out)
}

// Sends a hello greeting
func (s *_gin_GreeterService) _gin_handler_SayHello_1(ctx *gin.Context) {
	var in SayHelloRequest
	if err := ctx.ShouldBindUri(&in); err != nil {
		ctx.Error(err)
		return
	}
	if err := ctx.ShouldBind(&in); err != nil {
		ctx.Error(err)
		return
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request.Context(), md)
	out, err := s.server.(GreeterServiceGinServer).SayHello(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, out)
}

// Sends a hi greeting
func (s *_gin_GreeterService) _gin_handler_SayHi_0(ctx *gin.Context) {
	var in SayHiRequest
	if err := ctx.ShouldBindUri(&in); err != nil {
		ctx.Error(err)
		return
	}
	if err := ctx.ShouldBind(&in); err != nil {
		ctx.Error(err)
		return
	}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx.Request.Context(), md)
	out, err := s.server.(GreeterServiceGinServer).SayHi(newCtx, &in)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, out)
}

func (s *_gin_GreeterService) registerService() {

	// Sends a hello greeting
	s.router.Handle("GET", "/v1/helloworld.Greeter/SayHello/:name", s._gin_handler_SayHello_0)

	// Sends a hello greeting
	s.router.Handle("POST", "/v1/helloworld.Greeter/SayHello", s._gin_handler_SayHello_1)

	// Sends a hi greeting
	s.router.Handle("POST", "/helloworld.v1.GreeterService/SayHi", s._gin_handler_SayHi_0)

}