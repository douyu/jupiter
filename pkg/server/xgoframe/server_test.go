package xgoframe

import (
	"github.com/douyu/jupiter"
	"github.com/gogf/gf/net/ghttp"
	"testing"
)

func TestServer_Serve(t *testing.T) {
	t.Log("test over")
	var app jupiter.Application

	_ = app.Startup()
	_ = app.Serve(startGfServer())
	_ = app.Run()
}

func startGfServer() *Server  {
	serve := DefaultConfig().Build()

	serve.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello xgoframe for jupiter！")
	})


	return serve
}
