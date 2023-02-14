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

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	Name           string // config's name
	BalancerName   string
	Addr           string
	DialTimeout    time.Duration
	ReadTimeout    time.Duration
	KeepAlive      *keepalive.ClientParameters
	RegistryConfig string

	logger      *xlog.Logger
	dialOptions []grpc.DialOption

	SlowThreshold time.Duration

	Debug                      bool
	DisableSentinelInterceptor bool
	DisableTraceInterceptor    bool
	DisableAidInterceptor      bool
	DisableTimeoutInterceptor  bool
	DisableMetricInterceptor   bool
	DisableAccessInterceptor   bool
	AccessInterceptorLevel     string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            cast.ToDuration("3s"),
		ReadTimeout:            cast.ToDuration("1s"),
		SlowThreshold:          cast.ToDuration("600ms"),
		AccessInterceptorLevel: "info",
		KeepAlive: &keepalive.ClientParameters{
			Time:    5 * time.Minute,
			Timeout: 20 * time.Second,
		},
		RegistryConfig: constant.ConfigKey("registry.default"),
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("grpc." + name))
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config.Name = key
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
func (config *Config) Build() (*grpc.ClientConn, error) {
	config.logger = xlog.Jupiter().Named(ecode.ModClientGrpc)

	if config.Debug {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(config.Addr)),
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
			grpc.WithChainUnaryInterceptor(TraceUnaryClientInterceptor()),
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

	if !config.DisableSentinelInterceptor {
		config.dialOptions = append(config.dialOptions,
			grpc.WithChainUnaryInterceptor(sentinelUnaryClientInterceptor(config.Addr)),
		)
	}

	return newGRPCClient(config)
}

// Singleton returns a singleton client conn.
func (config *Config) Singleton() (*grpc.ClientConn, error) {
	if val, ok := singleton.Load(constant.ModuleClientGrpc, config.Name); ok && val != nil {
		return val.(*grpc.ClientConn), nil
	}

	cc, err := config.Build()
	if err != nil {
		xlog.Jupiter().Error("build grpc client failed", zap.Error(err))
		return nil, err
	}

	singleton.Store(constant.ModuleClientGrpc, config.Name, cc)

	return cc, nil
}

// MustBuild panics when error found.
func (config *Config) MustBuild() *grpc.ClientConn {
	return lo.Must(config.Build())
}

// MustSingleton panics when error found.
func (config *Config) MustSingleton() *grpc.ClientConn {
	return lo.Must(config.Singleton())
}
