package server

import (
	"github.com/douyu/jupiter"
	"github.com/google/wire"
	"uuid/internal/app/uuidserver/controller"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(controller.Options), "*"),
	NewGrpcServer,
	NewHttpServer,
)

type Options struct {
	http *HttpServer
	grpc *GrpcServer
}

func initApp(app *jupiter.Application, opts Options) error {
	// http
	if err := app.Serve(opts.http); err != nil {
		return err
	}

	// grpc
	if err := app.Serve(opts.grpc); err != nil {
		return err
	}

	return nil
}
