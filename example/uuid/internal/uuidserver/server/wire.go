//go:build wireinject
// +build wireinject

package server

import (
	"github.com/douyu/jupiter"
	"github.com/google/wire"
	"uuid/internal/uuidserver/controller"
	"uuid/internal/uuidserver/service"
)

func InitApp(app *jupiter.Application) error {
	panic(wire.Build(
		wire.Struct(new(Options), "*"),
		controller.ProviderSet,
		service.ProviderSet,
		ProviderSet,
		initApp,
	))
}
