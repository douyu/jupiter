package application

import (
	"os"

	"github.com/douyu/jupiter/pkg/flag"
)

var (
	appHost = os.Getenv("JUPITER_APP_HOST")
)

// EnvServerHost gets JUPITER_APP_HOST.
func EnvServerHost() string {
	host := flag.String("host")
	if host != "" {
		return host
	}

	if appHost == "" {
		return "127.0.0.1"
	}
	return appHost
}
