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

package redis

import (
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
)

// RedisClusterConfig redis集群配置信息
type RedisClusterConfig struct {
	// Addrs 集群实例配置地址
	Addrs []string `json:"addrs"`
	// Password 密码
	Password string `json:"password"`
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int `json:"db"`
	// PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	PoolSize int `json:"poolSize"`
	// MaxRedirects 网络相关的错误最大重试次数 默认8次
	MaxRedirects int `json:"maxRedirects"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `json:"minIdleConns"`
	// DialTimeout 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout"`
	// ReadTimeout 读超时 默认3s
	ReadTimeout time.Duration `json:"readTimeout"`
	// WriteTimeout 读超时 默认3s
	WriteTimeout time.Duration `json:"writeTimeout"`
	// IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	IdleTimeout time.Duration `json:"idleTimeout"`
	// ReadOnly 集群模式 在从属节点上启用读模式
	ReadOnly bool `json:"readOnly"`
	// Debug开关
	Debug bool `json:"debug"`
	// 是否开启链路追踪，开启以后。使用DoCotext的请求会被trace
	EnableTrace bool `json:"enableTrace"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold"`
	// OnDialError panic|error
	OnDialError string `json:"level"`
	logger      *xlog.Logger
}

// RedisConfig 单节点redis配置
type RedisConfig struct {
	// Addr 节点连接地址
	Addr string `json:"addr"`
	// Password 密码
	Password string `json:"password"`
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int `json:"db"`
	// PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	PoolSize int `json:"poolSize"`
	// MaxRedirects 网络相关的错误最大重试次数 默认8次
	MaxRetries int `json:"maxRetries"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `json:"minIdleConns"`
	// DialTimeout 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout"`
	// ReadTimeout 读超时 默认3s
	ReadTimeout time.Duration `json:"readTimeout"`
	// WriteTimeout 读超时 默认3s
	WriteTimeout time.Duration `json:"writeTimeout"`
	// IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
	IdleTimeout time.Duration `json:"idleTimeout"`
	// Debug开关
	Debug bool `json:"debug"`
	// 是否开启链路追踪，开启以后。使用DoCotext的请求会被trace
	EnableTrace bool `json:"enableTrace"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold"`
	// OnDialError panic|error
	OnDialError string `json:"level"`
	logger      *xlog.Logger
}

// StdRedisConfig ...
func StdRedisConfig(name string) RedisConfig {
	return RawRedisConfig("jupiter.redis." + name + ".stub")
}

// RawRedisConfig ...
func RawRedisConfig(key string) RedisConfig {
	var config = DefaultRedisConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal config", xlog.String("key", key), xlog.Any("config", config))
	}
	return config
}

// Build ...
func (config RedisConfig) Build() *RedisStub {
	if config.Addr == "" {
		config.logger.Panic("addr empty stub config", xlog.Any("config", config))
	}
	return newRedisStub(&config)
}

// StdRedisClusterConfig ...
func StdRedisClusterConfig(name string) RedisClusterConfig {
	return RawRedisClusterConfig("jupiter.redis." + name + ".cluster")
}

// RawRedisConfig ...
func RawRedisClusterConfig(key string) RedisClusterConfig {
	var config = DefaultRedisClusterConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal config", xlog.String("key", key), xlog.Any("config", config))

	}
	return config
}

// Build ...
func (config RedisClusterConfig) Build() *RedisClusterStub {
	if len(config.Addrs) == 0 {
		config.logger.Panic("cluster addr empty stub config")
	}
	return newRedisClusterStub(&config)
}

// DefaultRedisConfig ...
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		DB:            0,
		PoolSize:      10,
		MaxRetries:    3,
		MinIdleConns:  100,
		DialTimeout:   xtime.Duration("1s"),
		ReadTimeout:   xtime.Duration("1s"),
		WriteTimeout:  xtime.Duration("1s"),
		IdleTimeout:   xtime.Duration("60s"),
		Debug:         false,
		EnableTrace:   false,
		SlowThreshold: xtime.Duration("250ms"),
		OnDialError:   "panic",
		logger:        xlog.DefaultLogger,
	}
}

// DefaultRedisClusterConfig ...
func DefaultRedisClusterConfig() RedisClusterConfig {
	return RedisClusterConfig{
		DB:            0,
		PoolSize:      10,
		MaxRedirects:  3,
		MinIdleConns:  100,
		DialTimeout:   xtime.Duration("1s"),
		ReadTimeout:   xtime.Duration("1s"),
		WriteTimeout:  xtime.Duration("1s"),
		IdleTimeout:   xtime.Duration("60s"),
		ReadOnly:      false,
		Debug:         false,
		EnableTrace:   false,
		SlowThreshold: xtime.Duration("250ms"),
		OnDialError:   "panic",
		logger:        xlog.DefaultLogger,
	}
}
