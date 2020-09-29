// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xgoframe

import (
	"context"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/xlog"
	//"github.com/douyu/jupiter/pkg/ecode"
	//"github.com/douyu/jupiter/pkg/xlog"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

//Server is server core struct
type Server struct {
	*ghttp.Server
	config *Config
}

func newServer(config *Config) *Server {
	s := new(Server)
	serve := g.Server()
	serve.SetAddr(config.Address())

	s.Server = serve
	s.config = config

	return s
}

//Serve ..
func (s *Server) Serve() error {
	routes := s.GetRouterArray()

	for i := 0; i < len(routes); i++ {
		s.config.logger.Info("add route ", xlog.FieldMethod(routes[i].Method), xlog.String("path", routes[i].Route))
	}
	s.Run()

	return nil
}

//Stop ..
func (s *Server) Stop() error {
	return s.Shutdown()
}

//GracefulStop ..
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Stop()
}

//Info ..
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.config.Address()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}
