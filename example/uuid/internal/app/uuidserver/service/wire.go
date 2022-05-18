//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
	"uuid/internal/pkg/redis"
)

func createMockUuidService() *Uuid {
	panic(wire.Build(
		NewUuidService,
		redis.ProviderSet,
		wire.Struct(new(Options), "*"),
	))
}
