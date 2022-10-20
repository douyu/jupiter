// Copyright 2022 Douyu
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

package gorm

import (
	"time"

	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// StdConfig 标准配置，规范配置文件头
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("mysql." + name))
}

// RawConfig 传入mapstructure格式的配置
func RawConfig(key string) *Config {
	config := DefaultConfig()
	config.Name = key

	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("unmarshal config", xlog.FieldErr(err), xlog.FieldName(key))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}

	return config
}

// Config options
type Config struct {
	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug开关
	Debug bool `json:"debug" toml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// 拨超时时间
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	// 自动使用影子表
	AutoShadowTable bool `json:"autoShadowTable" toml:"autoShadowTable"`
	raw             interface{}
	// 记录错误sql时,是否打印包含参数的完整sql语句
	// select * from aid = ?;
	// select * from aid = 288016;
	DetailSQL bool `json:"detailSql" toml:"detailSql"`
	// 重试次数
	Retry int `json:"retry" toml:"retry"`
	// 重试等待时间
	RetryWaitTime time.Duration `json:"retryWaitTime" toml:"retryWaitTime"`

	gormConfig gorm.Config `json:"-" toml:"-"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: cast.ToDuration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   cast.ToDuration("500ms"),
		DialTimeout:     cast.ToDuration("3s"),
		AutoShadowTable: false,
		raw:             nil,
		Retry:           2,
		RetryWaitTime:   cast.ToDuration("200ms"),
	}
}

func (config *Config) WithGormConfig(gormConfig gorm.Config) *Config {
	config.gormConfig = gormConfig

	return config
}

func (config *Config) MustBuild() *gorm.DB {
	return dial(config.Name, config)
}

func (config *Config) MustSingleton() *gorm.DB {
	if val, ok := singleton.Load(constant.ModuleStoreGorm, config.Name); ok && val != nil {
		return val.(*gorm.DB)
	}

	db := dial(config.Name, config)

	singleton.Store(constant.ModuleStoreGorm, config.Name, db)
	return db
}
