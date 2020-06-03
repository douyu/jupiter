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

package grpc

import (
	"time"

	"github.com/douyu/jupiter/pkg/ecode"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	BalancerName string
	Address      string
	Block        bool
	DialTimeout  time.Duration
	Direct       bool
	OnDialError  string // panic | error
	KeepAlive    *keepalive.ClientParameters
	logger       *xlog.Logger
	dialOptions  []grpc.DialOption
	// resolver     resolver.Builder

	Debug         bool
	DisableTrace  bool
	DisableMetric bool
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		logger:       xlog.JupiterLogger.With(xlog.FieldMod(ecode.ModClientGrpc)),
		BalancerName: roundrobin.Name, // roundrobin by default
		DialTimeout:  time.Second * 3,
		OnDialError:  "panic",
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("jupiter.client." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("client grpc parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// WithDialOption ...
func (config *Config) WithDialOption(opts ...grpc.DialOption) *Config {
	if config.dialOptions == nil {
		config.dialOptions = make([]grpc.DialOption, 0)
	}
	config.dialOptions = append(config.dialOptions, opts...)
	return config
}

// Build ...
func (config *Config) Build() *grpc.ClientConn {
	if config.Debug {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(config.Address)),
		)
	}
	if !config.DisableTrace {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	client := newGRPCClient(*config)
	return client
}
