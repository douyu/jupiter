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

package xecho

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Server ...
type Server struct {
	*echo.Echo
	config   *Config
	listener net.Listener
	// registerer registry.Registry
}

func newServer(config *Config) (*Server, error) {
	var (
		listener net.Listener
		err      error
	)

	if config.EnableTLS {
		var cert, key []byte
		cert, err = ioutil.ReadFile(config.CertFile)
		if err != nil {
			return nil, errors.Wrap(err, "read cert failed")
		}

		key, err = ioutil.ReadFile(config.PrivateFile)
		if err != nil {
			return nil, errors.Wrap(err, "read private failed")
		}

		tlsConfig := new(tls.Config)
		tlsConfig.Certificates = make([]tls.Certificate, 1)

		if tlsConfig.Certificates[0], err = tls.X509KeyPair(cert, key); err != nil {
			return nil, errors.Wrap(err, "X509KeyPair failed")
		}
		listener, err = tls.Listen("tcp", config.Address(), tlsConfig)
	} else {
		listener, err = net.Listen("tcp", config.Address())
	}
	if err != nil {
		// config.logger.Panic("new xecho server err", xlog.FieldErrKind(ecode.ErrKindListenErr), xlog.FieldErr(err))
		return nil, errors.Wrapf(err, "create xecho server failed")
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port
	return &Server{
		Echo:     echo.New(),
		config:   config,
		listener: listener,
	}, nil
}

func (s *Server) Healthz() bool {
	if s.Echo.Listener == nil {
		return false
	}

	conn, err := s.Echo.Listener.Accept()
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// Server implements server.Server interface.
func (s *Server) Serve() error {
	s.Echo.Logger.SetOutput(os.Stdout)
	s.Echo.Debug = s.config.Debug
	s.Echo.HideBanner = true
	s.Echo.StdLogger = zap.NewStdLog(xlog.Jupiter())
	for _, route := range s.Echo.Routes() {
		s.config.logger.Info("add route", xlog.FieldMethod(route.Method), xlog.String("path", route.Path))
	}

	var err error

	if s.config.EnableTLS {
		s.Echo.TLSListener = s.listener
		err = s.Echo.StartTLS("", s.config.CertFile, s.config.PrivateFile)
	} else {
		s.Echo.Listener = s.listener
		err = s.Echo.Start("")
	}

	if err != http.ErrServerClosed {
		return err
	}
	s.config.logger.Info("close echo", xlog.FieldAddr(s.config.Address()))
	return nil
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	return s.Echo.Close()
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.config.ServiceAddress != "" {
		serviceAddr = s.config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceProvider),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
