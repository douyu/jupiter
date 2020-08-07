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

	"github.com/go-redis/redis"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
)

const (
	//ClusterMode using clusterClient
	ClusterMode string = "cluster"
	//StubMode using reidsClient
	StubMode string = "stub"
)

// Config for redis, contains RedisStubConfig and RedisClusterConfig
type Config struct {
	// Addrs 实例配置地址
	Addrs []string `json:"addrs"`
	// Addr stubConfig 实例配置地址
	Addr string `json:"addr"`
	// Mode Redis模式 cluster|stub
	Mode string `json:"mode"`
	// Password 密码
	Password string `json:"password"`
	// DB，默认为0, 一般应用不推荐使用DB分片
	DB int `json:"db"`
	// PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	PoolSize int `json:"poolSize"`
	// MaxRetries 网络相关的错误最大重试次数 默认8次
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
	// ReadOnly 集群模式 在从属节点上启用读模式
	ReadOnly bool `json:"readOnly"`
	// 是否开启链路追踪，开启以后。使用DoCotext的请求会被trace
	EnableTrace bool `json:"enableTrace"`
	// 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowThreshold time.Duration `json:"slowThreshold"`
	// OnDialError panic|error
	OnDialError string `json:"level"`
	logger      *xlog.Logger
}

// DefaultRedisConfig default config ...
func DefaultRedisConfig() Config {
	return Config{
		DB:            0,
		PoolSize:      10,
		MaxRetries:    3,
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
		logger:        xlog.JupiterLogger,
	}
}

// StdRedisConfig ...
func StdRedisConfig(name string) Config {
	return RawRedisConfig("jupiter.redis." + name)
}

// RawRedisConfig ...
func RawRedisConfig(key string) Config {
	var config = DefaultRedisConfig()

	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal redisConfig",
			xlog.String("key", key),
			xlog.Any("redisConfig", config),
			xlog.String("error", err.Error()))
	}
	return config
}

// Build ...
func (config Config) Build() *Redis {
	count := len(config.Addrs)
	if count < 1 {
		config.logger.Panic("no address in redis config", xlog.Any("config", config))
	}
	if len(config.Mode) == 0 {
		config.Mode = StubMode
		if count > 1 {
			config.Mode = ClusterMode
		}
	}
	var client redis.Cmdable
	switch config.Mode {
	case ClusterMode:
		if count == 1 {
			config.logger.Warn("redis config has only 1 address but with cluster mode")
		}
		client = config.buildCluster()
	case StubMode:
		if count > 1 {
			config.logger.Warn("redis config has more than 1 address but with stub mode")
		}
		client = config.buildStub()
	default:
		config.logger.Panic("redis mode must be one of (stub, cluster)")
	}
	return &Redis{
		Config: &config,
		Client: client,
	}
}

func (config Config) buildStub() *redis.Client {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         config.Addrs[0],
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})

	if err := stubClient.Ping().Err(); err != nil {
		switch config.OnDialError {
		case "panic":
			config.logger.Panic("dial redis fail", xlog.Any("err", err), xlog.Any("config", config))
		default:
			config.logger.Error("dial redis fail", xlog.Any("err", err), xlog.Any("config", config))
		}
	}

	return stubClient

}

func (config Config) buildCluster() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		MaxRedirects: config.MaxRetries,
		ReadOnly:     config.ReadOnly,
		Password:     config.Password,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	if err := clusterClient.Ping().Err(); err != nil {
		switch config.OnDialError {
		case "panic":
			config.logger.Panic("start cluster redis", xlog.Any("err", err))
		default:
			config.logger.Error("start cluster redis", xlog.Any("err", err))
		}
	}
	return clusterClient
}

// StdRedisStubConfig ...
func StdRedisStubConfig(name string) Config {
	return RawRedisStubConfig("jupiter.redis." + name + ".stub")
}

// RawRedisStubConfig ...
func RawRedisStubConfig(key string) Config {
	var config = DefaultRedisConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("unmarshal config",
			xlog.String("key", key),
			xlog.Any("config", config),
			xlog.Any("error", err))
	}
	config.Addrs = []string{config.Addr}
	config.Mode = StubMode
	return config
}

// StdRedisClusterConfig ...
func StdRedisClusterConfig(name string) Config {
	return RawRedisClusterConfig("jupiter.redis." + name + ".cluster")
}

// RawRedisClusterConfig ...
func RawRedisClusterConfig(key string) Config {
	var config = DefaultRedisConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("unmarshal config",
			xlog.String("key", key),
			xlog.Any("config", config),
			xlog.Any("error", err))
	}
	config.Mode = ClusterMode
	return config
}
