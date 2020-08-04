// @author ccb1900
// @url https://github.com/ccb1900
// @description this a server wrapper for goframe with jupiter
// it is necessary to implement Server interface
// You should add middleware for metrics,trace and recovery
package xgoframe

import (
	"context"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/xlog"
	"strconv"

	//"github.com/douyu/jupiter/pkg/ecode"
	//"github.com/douyu/jupiter/pkg/xlog"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type Server struct {
	*ghttp.Server
	config *Config
}

func newServer(config *Config) *Server {
	s := new(Server)
	serve := g.Server()
	serve.SetPort(config.Port)

	s.Server = serve
	s.config = config

	return s
}

func (s *Server)Serve() error  {
	routes := s.GetRouterArray()

	for i := 0; i < len(routes); i++ {
		s.config.logger.Info("add route ",xlog.FieldMethod(routes[i].Method),xlog.String("path",routes[i].Route))
	}
	s.Run()

	return nil
}

func (s *Server)Stop() error  {
	return s.Shutdown()
}

func (s *Server)GracefulStop(ctx context.Context) error  {
	return s.Stop()
}

func (s *Server)Info() *server.ServiceInfo  {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.config.Host +":"+ strconv.Itoa(s.config.Port)),
		server.WithKind(constant.ServiceProvider),
	)
	info.Name = info.Name + "." + ModName
	return &info
}
