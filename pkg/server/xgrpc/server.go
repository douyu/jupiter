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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// Server ...
type Server struct {
	*grpc.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) (*Server, error) {
	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{defaultStreamServerInterceptor(config.logger, config)},
		config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor(config.logger, config)},
		config.unaryInterceptors...,
	)

	if !config.DisableTrace {
		unaryInterceptors = append(
			[]grpc.UnaryServerInterceptor{NewTraceUnaryServerInterceptor()},
			unaryInterceptors...,
		)

		streamInterceptors = append(
			[]grpc.StreamServerInterceptor{NewTraceStreamServerInterceptor()},
			streamInterceptors...,
		)
	}

	if !config.DisableMetric {
		unaryInterceptors = append(
			[]grpc.UnaryServerInterceptor{prometheusUnaryServerInterceptor},
			unaryInterceptors...,
		)

		streamInterceptors = append(
			[]grpc.StreamServerInterceptor{prometheusStreamServerInterceptor},
			streamInterceptors...,
		)
	}

	if !config.DisableSentinel {
		unaryInterceptors = append(
			[]grpc.UnaryServerInterceptor{NewSentinelUnaryServerInterceptor()},
			unaryInterceptors...,
		)
	}

	if config.EnableTLS {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.PrivateFile)
		if err != nil {
			return nil, errors.Wrap(err, "tls.LoadX509KeyPair failed")
		}

		certPool := x509.NewCertPool()
		rootBuf, err := ioutil.ReadFile(config.CaFile)
		if err != nil {
			return nil, errors.Wrap(err, "ioutil.ReadFile failed")
		}
		if !certPool.AppendCertsFromPEM(rootBuf) {
			return nil, errors.New("certPool.AppendCertsFromPEM failed")
		}

		tlsConf := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		}

		config.serverOptions = append(config.serverOptions,
			grpc.Creds(credentials.NewTLS(tlsConf)),
		)
	}

	config.serverOptions = append(config.serverOptions,
		grpc.StreamInterceptor(StreamInterceptorChain(streamInterceptors...)),
		grpc.UnaryInterceptor(UnaryInterceptorChain(unaryInterceptors...)),
	)

	newServer := grpc.NewServer(config.serverOptions...)
	listener, err := net.Listen(config.Network, config.Address())
	if err != nil {
		return nil, errors.Wrap(err, "net.Listen failed")
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	reflection.Register(newServer)

	return &Server{
		Server:   newServer,
		listener: listener,
		Config:   config,
	}, nil
}

func (s *Server) Healthz() bool {
	return true
}

// Server implements server.Server interface.
func (s *Server) Serve() error {
	// display grpc server method list
	for fm, info := range s.GetServiceInfo() {
		for _, method := range info.Methods {
			fmt.Printf("[GRPC] \x1b[34m%8s\x1b[0m.%s\n", fm, method.Name)
		}
	}
	// display grpc server addr
	fmt.Printf("[GRPC] \x1b[33m%8s\x1b[0m %s\n", "Listen On", s.listener.Addr().String())
	err := s.Server.Serve(s.listener)
	return err
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	s.Server.Stop()
	return nil
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	serviceAddress := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddress = s.Config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(serviceAddress),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}
