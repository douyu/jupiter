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

package xorm

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-xorm/xorm"
)

type Config struct {

	Name string
	// DSN地址: mysql://root:secret@tcp(127.0.0.1:3306)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug开关
	Debug bool `json:"debug" toml:"debug"`

	//打印sql
	ShowSql bool `json:"showSql" toml:"showSql"`

	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// 创建连接的错误级别，=panic时，如果创建失败，立即panic
	OnDialError string `json:"level" toml:"level"`

	raw          interface{}
	logger       *xlog.Logger
	dsnCfg       *DSN
}

func StdConfig(name string ) *Config {

	return RawConfig("jupiter.mysql." + name)
}

func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config, conf.TagName("toml")); err != nil {
		xlog.Panic("unmarshal key", xlog.FieldMod("xorm"), xlog.FieldErr(err), xlog.FieldKey(key))
	}
	config.Name = key
	return config
}

func DefaultConfig() *Config {

	return &Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		OnDialError:     "panic",
		raw:             nil,
		logger:          xlog.JupiterLogger,
	}
}


func (config *Config) Build () *xorm.Engine {
	var err error

	config.dsnCfg, err = ParseDSN(config.DSN)

	if err == nil {
		config.logger.Info(ecode.MsgClientMysqlOpenStart, xlog.FieldMod("xorm"), xlog.FieldAddr(config.dsnCfg.Addr), xlog.FieldName(config.dsnCfg.DBName))
	} else {
		config.logger.Panic(ecode.MsgClientMysqlOpenStart, xlog.FieldMod("xorm"), xlog.FieldErr(err))
	}

	eng , err := Open("mysql", config)

	if err != nil {

		if config.OnDialError == "panic" {

			config.logger.Panic("open mysql", xlog.FieldMod("xorm"), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldAddr(config.dsnCfg.Addr), xlog.FieldValueAny(config))
		}else{

			metric.LibHandleCounter.Inc(metric.TypeXorm, config.Name+".ping", config.dsnCfg.Addr, "open err")
			config.logger.Error("open mysql", xlog.FieldMod("xorm"), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldAddr(config.dsnCfg.Addr), xlog.FieldValueAny(config))
			return eng
		}

	}

	if err = eng.Ping(); err != nil {
		config.logger.Panic("ping mysql", xlog.FieldMod("xorm"), xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}

	return eng
}