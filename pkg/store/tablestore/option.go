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

package tstore

import (
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config ...
type (
	Config struct {
		Name string
		// Debug 开关
		Debug bool `toml:"debug" json:"debug"`
		// 指标采集开关
		EnableMetric bool `toml:"enableMetric" json:"enableMetric"`
		// 连接端点
		EndPoint string `toml:"endPoint" json:"endPoint"`
		// 实例
		Instance string `toml:"instance" json:"instance"`
		// accessKeyId
		AccessKeyId string `toml:"accessKeyId" json:"accessKeyId"`
		// accessKeySecret
		AccessKeySecret string `toml:"accessKeySecret" json:"accessKeySecret"`
		// 安全密钥
		SecurityToken string `toml:"securityToken" json:"securityToken"`
		// 重试次数
		RetryTimes uint `toml:"retryTimes" json:"retryTimes"`
		// 慢日志阈值
		SlowThreshold time.Duration `toml:"slowThreshold" json:"slowThreshold"`
		// 最大重试时间
		MaxRetryTime time.Duration `toml:"maxRetryTime" json:"maxRetryTime"`
		// 连接超时时间
		ConnectionTimeout time.Duration `toml:"connectionTimeout" json:"connectionTimeout"`
		// 请求超时时间
		RequestTimeout time.Duration `toml:"requestTimeout" json:"requestTimeout"`
		// 最大空闲连接数
		MaxIdleConnections int `toml:"maxIdleConnections" json:"maxIdleConnections"`
		// 访问日志开关
		EnableAccessLog bool `toml:"enableAccessLog" json:"enableAccessLog"`
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("tablestore." + name))
}

// RawConfig jupiter.tablestore.demo
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config.Name = key
	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("tablestore unmarshal config", xlog.FieldErr(err), xlog.FieldName(key))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}
	return config
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Debug:              false,
		EnableMetric:       true,
		EnableAccessLog:    false,
		EndPoint:           "",
		Instance:           "",
		AccessKeyId:        "",
		AccessKeySecret:    "",
		SecurityToken:      "",
		RetryTimes:         1,
		SlowThreshold:      time.Second * 1,
		MaxRetryTime:       time.Second * 5,
		ConnectionTimeout:  time.Second * 15,
		RequestTimeout:     time.Second * 30,
		MaxIdleConnections: 2000,
	}
}
func (config *Config) MustSingleton() *tablestore.TableStoreClient {
	if val, ok := singleton.Load(constant.ModuleStoreTableStore, config.Name); ok && val != nil {
		return val.(*tablestore.TableStoreClient)
	}
	ts := newTs(config)
	singleton.Store(constant.ModuleStoreTableStore, config.EndPoint, ts)
	return ts
}
func (config *Config) MustBuild() *tablestore.TableStoreClient {
	return newTs(config)
}
