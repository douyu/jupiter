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

package mongodb

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// StdConfig 配置
func StdConfig(name string) *Config {
	return RawConfig("jupiter.mongodb." + name)
}

// RawConfig 传入mapstructure格式的配置
// 例如：jupiter.mongodb.config
func RawConfig(key string) *Config {
	var cfg = DefaultConfig()
	if err := conf.UnmarshalKey(key, cfg, conf.TagName("toml")); err != nil {
		xlog.Panic("unmarshal key", xlog.FieldMod("mongodb"), xlog.FieldErr(err), xlog.FieldKey(key))
	}
	return cfg
}

// 配置
type Config struct {
	// mongodb uri 链接地址
	// mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	URI             string        `json:"uri" toml:"uri"`
	Debug           bool          `json:"debug" toml:"debug"`                           // debug 模式
	ConnectTimeout  time.Duration `json:"connect_timeout" toml:"connect_timeout"`       // 连接超时时间
	MaxConnIdleTime time.Duration `json:"max_conn_idle_time" toml:"max_conn_idle_time"` // 最大空闲连接时间
	MaxPoolSize     uint64        `json:"max_pool_size" toml:"max_pool_size"`           // 连接池最大连接数
	MinPoolSize     uint64        `json:"min_pool_size" toml:"min_pool_size"`           // 连接池最小连接数
	DisableMetric   bool          `json:"disable_metric" toml:"disable_metric"`         // 关闭指标采集
	DisableTrace    bool          `json:"disable_trace" toml:"disable_trace"`           // 关闭链路追踪
	logger          *xlog.Logger  // logger
}

// 默认配置

func DefaultConfig() *Config {
	return &Config{
		URI:             "", // 空的uri
		Debug:           false,
		MaxConnIdleTime: xtime.Duration("60s"), //  60s
		MaxPoolSize:     100,                   // 100个链接
		MinPoolSize:     10,                    // 10个链接
		ConnectTimeout:  xtime.Duration("10s"), // 10s
		DisableMetric:   false,
		DisableTrace:    false,
		logger:          xlog.DefaultLogger,
	}
}

// WithLogger ...
func (config *Config) WithLogger(log *xlog.Logger) *Config {
	config.logger = log
	return config
}

// build connect mongodb
func (config *Config) Build() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI))
	if err != nil {
		config.logger.Panic("connect mongodb", xlog.FieldMod("mongodb"), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}
	return client
}
