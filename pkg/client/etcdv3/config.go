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

package etcdv3

import (
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Config ...
type (
	Config struct {
		Name      string   `json:"name"`
		Endpoints []string `json:"endpoints"`
		CertFile  string   `json:"certFile"`
		KeyFile   string   `json:"keyFile"`
		CaCert    string   `json:"caCert"`
		BasicAuth bool     `json:"basicAuth"`
		UserName  string   `json:"userName"`
		Password  string   `json:"-"`
		// 连接超时时间
		ConnectTimeout time.Duration `json:"connectTimeout"`
		Secure         bool          `json:"secure"`
		// 自动同步member list的间隔
		AutoSyncInterval time.Duration `json:"autoAsyncInterval"`
		TTL              int           // 单位：s
		EnableTrace      bool          `json:"enableTrace" toml:"enableTrace"`

		logger *xlog.Logger
	}
)

func (config *Config) BindFlags(fs *flag.FlagSet) {
	fs.BoolVar(&config.Secure, "insecure-etcd", true, "--insecure-etcd=true")
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Endpoints:      []string{"http://localhost:2379"},
		BasicAuth:      false,
		ConnectTimeout: cast.ToDuration("5s"),
		Secure:         false,
		EnableTrace:    true,
		logger:         xlog.Jupiter().With(xlog.FieldMod("client.etcd")),
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("etcdv3." + name))
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config.Name = key

	if err := conf.UnmarshalKey(key, config); err != nil {
		config.logger.Panic("client etcd parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}

	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// Build ...
func (config *Config) Build() (*Client, error) {
	return newClient(config)
}

// Singleton returns a singleton client conn.
func (config *Config) Singleton() (*Client, error) {
	if client, ok := singleton.Load(constant.ModuleRegistryEtcd, config.Name); ok && client != nil {
		return client.(*Client), nil
	}

	client, err := config.Build()
	if err != nil {
		xlog.Jupiter().Error("build etcd client failed", zap.Error(err))
		return nil, err
	}

	singleton.Store(constant.ModuleRegistryEtcd, config.Name, client)

	return client, nil
}

// MustBuild panics when error found.
func (config *Config) MustBuild() *Client {
	return lo.Must(config.Build())
}

// MustSingleton panics when error found.
func (config *Config) MustSingleton() *Client {
	return lo.Must(config.Singleton())
}
