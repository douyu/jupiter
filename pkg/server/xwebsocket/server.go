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

package xgrpc

import (
	"context"
	"net"
	"net/http"

	"github.com/douyu/jupiter/pkg"
	"github.com/gorilla/websocket"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Server ...
type Server struct {
	config   *Config
	listener net.Listener
	upgrader *websocket.Upgrader
}

func newServer(config *Config) *Server {
	upgrader := &websocket.Upgrader{}
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		config.logger.Panic("new websockt server err", xlog.FieldErrKind(ecode.ErrKindListenErr), xlog.FieldErr(err))
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	return &Server{upgrader: upgrader, listener: listener, config: config}
}

// Serve implements server.Serve interface.
func (s *Server) Serve() error {
	s.Server = &http.Server{}
	err := http.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

//Upgrade get upgrage request
func (s *Server) Upgrade(uri string, fn func(*websocket.Conn, error)) error {
	http.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		ctx, err := s.upgrader.Upgrade(w, r, nil)
		if err == nil {
			defer ctx.Close()
		}
		fn(ctx, err)
	})

}

// Stop implements server.Stop interface
// it will terminate websocket server immediately
func (s *Server) Stop() error {
	s.listener.Close()
	return nil
}

// GracefulStop implements server.GracefulStop interface
// it will stop websocket server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.listener.Close()
	return nil
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:      pkg.Name(),
		Scheme:    "http",
		IP:        s.config.Host,
		Port:      s.config.Port,
		Weight:    0.0,
		Enable:    false,
		Healthy:   false,
		Metadata:  map[string]string{},
		Region:    "",
		Zone:      "",
		GroupName: "",
	}
}
