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
	"github.com/douyu/jupiter/pkg/util/xtime"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	Name         string // config's name
	BalancerName string
	Address      string
	Block        bool
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	Direct       bool
	OnDialError  string // panic | error
	KeepAlive    *keepalive.ClientParameters
	logger       *xlog.Logger
	dialOptions  []grpc.DialOption

	SlowThreshold time.Duration

	Debug                     bool
	DisableTraceInterceptor   bool
	DisableAidInterceptor     bool
	DisableTimeoutInterceptor bool
	DisableMetricInterceptor  bool
	DisableAccessInterceptor  bool
	AccessInterceptorLevel    string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		logger:                 xlog.JupiterLogger.With(xlog.FieldMod(ecode.ModClientGrpc)),
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 3,
		ReadTimeout:            xtime.Duration("1s"),
		SlowThreshold:          xtime.Duration("600ms"),
		OnDialError:            "panic",
		AccessInterceptorLevel: "info",
		Block:                  true,
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

	if !config.DisableAidInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(aidUnaryClientInterceptor()),
		)
	}

	if !config.DisableTimeoutInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(config.logger, config.ReadTimeout, config.SlowThreshold)),
		)
	}

	if !config.DisableTraceInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	if !config.DisableAccessInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(config.logger, config.Name, config.AccessInterceptorLevel)),
		)
	}

	if !config.DisableMetricInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(metricUnaryClientInterceptor(config.Name)),
		)
	}

	return newGRPCClient(config)
}
