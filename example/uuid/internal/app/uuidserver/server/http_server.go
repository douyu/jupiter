package server

import (
	"github.com/douyu/jupiter/pkg/server/xecho"
	echo "github.com/labstack/echo/v4"
	"uuid/internal/app/uuidserver/controller"
)

type HttpServer struct {
	*xecho.Server
	controller.Options
}

func NewHttpServer(opts controller.Options) *HttpServer {
	s := xecho.StdConfig("http").MustBuild()

	s.GET("/snowflake_uuid", func(c echo.Context) error {
		return opts.UuidHTTP.GetUuidBySnowflake(c)
	})

	s.GET("/google_uuid_v4", func(c echo.Context) error {
		return opts.UuidHTTP.GetUuidByGoogleUUIDV4(c)
	})

	return &HttpServer{
		Server: s,
	}
}
