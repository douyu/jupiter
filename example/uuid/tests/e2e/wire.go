//go:build wireinject
// +build wireinject

package e2e

import (
	"uuid/internal/app/uuidserver/service"
	// "uuid/internal/pkg/mysql"
	// "uuid/internal/pkg/redis"
	"github.com/google/wire"
)

func CreateUuidService() *service.Uuid {
	panic(wire.Build(
		service.NewUuidService,
		// grpc.ProviderSet,
		wire.Struct(new(service.Options), "*"),
	))
}
