package xnet

import (
	"fmt"
	"net"

	"github.com/samber/lo"
)

// Address means the address of the service to be registered
func Address(listener net.Listener) string {
	host, port := lo.Must2(net.SplitHostPort(listener.Addr().String()))
	if host == "::" || host == "0.0.0.0" {
		host, _, _ = GetLocalMainIP()
	}

	return fmt.Sprintf("%s:%s", host, port)
}
