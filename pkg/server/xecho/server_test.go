package xecho

import (
	"github.com/douyu/jupiter"
	"github.com/labstack/echo/v4"
	"testing"
)

func TestServer_Serve(t *testing.T) {
	t.Log("test over")
	var app jupiter.Application

	_ = app.Startup()
	_ = app.Serve(startServer())
	_ = app.Run()
}

func startServer() *Server  {
	serve := DefaultConfig().Build()

	serve.GET("/", func(context echo.Context) error {
		return context.JSON(200,"test echo")
	})

	return serve
}
