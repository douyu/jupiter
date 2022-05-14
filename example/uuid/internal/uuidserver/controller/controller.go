package controller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUuidHTTPController)

type Options struct {
	UuidHTTP *UuidHTTP
}
