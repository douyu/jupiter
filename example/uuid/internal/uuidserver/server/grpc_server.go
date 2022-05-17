package server

import (
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	uuidv1 "uuid/gen/api/go/uuid/v1"
	"uuid/internal/uuidserver/service"
)

// var GrpcProviderSet = wire.NewSet(NewGrpcServer)

type GrpcServer struct {
	*xgrpc.Server
	Uuid *service.Uuid
}

func NewGrpcServer(opts *service.Uuid) *GrpcServer {
	return &GrpcServer{
		Server: xgrpc.StdConfig("grpc").MustBuild(),
		Uuid:   opts,
	}
}

func (s *GrpcServer) Mux() {
	uuidv1.RegisterUuidServer(s.Server.Server, s.Uuid)
}
