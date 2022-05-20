package server

import (
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	uuidv1 "uuid/gen/api/go/uuid/v1"
	"uuid/internal/app/uuidserver/controller"
)

// var GrpcProviderSet = wire.NewSet(NewGrpcServer)

type GrpcServer struct {
	*xgrpc.Server
	controller.Options
}

func NewGrpcServer(opts controller.Options) *GrpcServer {
	server := xgrpc.StdConfig("grpc").MustBuild()
	uuidv1.RegisterUuidServer(server.Server, opts.UuidGrpc)
	return &GrpcServer{
		Server: server,
	}
}
