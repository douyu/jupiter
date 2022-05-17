//go:build wireinject
// +build wireinject

package service

import (
	"github.com/google/wire"
)

func createMockUuidService() *Uuid {
	panic(wire.Build(
		NewUuidService,
		// grpc.ProviderSet,
		wire.Struct(new(Options), "*"),
	))
}
