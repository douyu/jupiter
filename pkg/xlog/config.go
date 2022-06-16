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

package xlog

import (
	"fmt"
	"log"
	"time"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init loggers")
		log.Printf("reload default logger with configKey: %s", ConfigEntry("default"))
		defaultLogger = RawConfig(constant.ConfigPrefix + ".logger.default").Build()

		log.Printf("reload default logger with configKey: %s", ConfigEntry("jupiter"))
		jupiterLogger = RawConfig(constant.ConfigPrefix + ".logger.jupiter").Build()
	})
}

var ConfigPrefix = constant.ConfigPrefix + ".logger"

// Config ...
type Config struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int
	MaxAge    int
	MaxBackup int
	// 日志磁盘刷盘间隔
	Interval      time.Duration
	CallerSkip    int
	Async         bool
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	Debug         bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
}

// Filename ...
func (config *Config) Filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

func ConfigEntry(name string) string {
	return ConfigPrefix + "." + name
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config, _ = conf.UnmarshalWithExpect(key, config).(*Config)
	config.configKey = key
	return config
}

// StdConfig Jupiter Standard logger config
func StdConfig(name string) *Config {
	return RawConfig(ConfigPrefix + "." + name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:          "jupiter_default.json",
		Dir:           pkg.LogDir(),
		Level:         "info",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    2,
		AddCaller:     true,
		Async:         true,
		Queue:         false,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: DefaultZapConfig(),
		Fields: []zap.Field{
			String("aid", pkg.AppID()),
			String("iid", pkg.AppInstance()),
		},
	}
}

// Build ...
func (config Config) Build() *Logger {
	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapConfig()
	}
	if config.Debug {
		config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	logger := newLogger(&config)

	return logger
}
