//go:build wireinject
// +build wireinject

package e2e

import (
	"github.com/google/wire"
	"uuid/internal/app/uuidserver/service"
	// "uuid/internal/pkg/mysql"
	"uuid/internal/pkg/redis"
)

func CreateUuidService() *service.Uuid {
	panic(wire.Build(
		service.NewUuidService,
		redis.ProviderSet,
		wire.Struct(new(service.Options), "*"),
	))
}
