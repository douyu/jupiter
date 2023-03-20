// Copyright 2022 Douyu
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

package xfasthttp

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xnet"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// Server ...
type Server struct {
	*fasthttp.Server
	config   *Config
	listener net.Listener
}

func newServer(config *Config) (*Server, error) {
	var (
		listener  net.Listener
		tlsConfig tls.Config
		err       error
	)

	if config.EnableTLS {
		cert, err := ioutil.ReadFile(config.CertFile)
		if err != nil {
			return nil, errors.Wrap(err, "read cert failed")
		}

		key, err := ioutil.ReadFile(config.PrivateFile)
		if err != nil {
			return nil, errors.Wrap(err, "read private failed")
		}

		tlsConfig.Certificates = make([]tls.Certificate, 1)

		if tlsConfig.Certificates[0], err = tls.X509KeyPair(cert, key); err != nil {
			return nil, errors.Wrap(err, "X509KeyPair failed")
		}

	}

	listener, err = net.Listen("tcp", config.Address())

	if err != nil {
		// config.logger.Panic("new fasthttp server err", xlog.FieldErrKind(ecode.ErrKindListenErr), xlog.FieldErr(err))
		return nil, errors.Wrapf(err, "create fasthttp server failed")
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	return &Server{
		Server: &fasthttp.Server{
			Concurrency:       config.Concurrency,
			ReadBufferSize:    config.ReadBufferSize,
			WriteBufferSize:   config.WriteBufferSize,
			ReduceMemoryUsage: config.ReduceMemoryUsage,
			TLSConfig:         &tlsConfig,
		},
		config:   config,
		listener: listener,
	}, nil
}

func (s *Server) Healthz() bool {
	return true
}

// Server implements server.Server interface.
func (s *Server) Serve() error {
	var err error

	s.Handler = recoverMiddleware(s.config)(s.Handler)

	if s.config.EnableTLS {
		err = s.Server.ServeTLS(s.listener, s.config.CertFile, s.config.PrivateFile)
	} else {
		err = s.Server.Serve(s.listener)
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
	return s.Server.Shutdown()
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown()
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(xnet.Address(s.listener)),
		server.WithKind(constant.ServiceProvider),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
