package controller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUuidHTTPController, NewUUuidGrpcController)

type Options struct {
	UuidHTTP *UuidHTTP
	UuidGrpc *UuidGrpc
}
