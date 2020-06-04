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
	"fmt"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"

	"github.com/douyu/jupiter/pkg/conf"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Host               string
	Port               int
	Network            string `json:"network" toml:"network"`
	DisableTrace       bool
	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor

	logger *xlog.Logger
}

// Jupiter Standard gRPC Server config
func StdConfig(name string) *Config {
	return RawConfig("jupiter.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("grpc server parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		serverOptions: []grpc.ServerOption{},
		Network:       "tcp4",
		Host:          "127.0.0.1",
		Port:          9092,
		logger:        xlog.JupiterLogger.With(xlog.FieldMod("server.grpc")),
	}
}

// WithServerOption ...
func (config *Config) WithServerOption(options ...grpc.ServerOption) Config {
	config.serverOptions = append(config.serverOptions, options...)
	return *config
}

// Build ...
func (config *Config) Build() *Server {
	config.streamInterceptors = []grpc.StreamServerInterceptor{
		config.RecoveryStreamServerInterceptor(),
		config.LoggerStreamServerIntercept(),
	}
	config.unaryInterceptors = []grpc.UnaryServerInterceptor{
		config.RecoveryUnaryServerInterceptor(),
		config.LoggerUnaryServerIntercept(),
	}
	if !config.DisableTrace {
		config.unaryInterceptors = append(config.unaryInterceptors, traceUnaryServerInterceptor)
		config.streamInterceptors = append(config.streamInterceptors, traceStreamServerInterceptor)
	}

	return newServer(config)
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// WithHost ...
func (config *Config) WithHost(host string) Config {
	config.Host = host
	return *config
}

// WithPort ...
func (config *Config) WithPort(port int) Config {
	config.Port = port
	return *config
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
