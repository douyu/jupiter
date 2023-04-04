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
	"os"
	"time"

	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Func ...
type (
	Field  = zap.Field
	Level  = zapcore.Level
	Logger = zap.Logger
	Option = zap.Option
)

var (
	// String ...
	String = zap.String
	// Any ...
	Any = zap.Any
	// Int64 ...
	Int64 = zap.Int64
	// Int ...
	Int = zap.Int
	// Int32 ...
	Int32 = zap.Int32
	// Uint ...
	Uint = zap.Uint
	// Duration ...
	Duration = zap.Duration
	// Durationp ...
	Durationp = zap.Durationp
	// Object ...
	Object = zap.Object
	// Namespace ...
	Namespace = zap.Namespace
	// Reflect ...
	Reflect = zap.Reflect
	// Skip ...
	Skip = zap.Skip()
	// ByteString ...
	ByteString = zap.ByteString
)

const (
	// defaultBufferSize sizes the buffer associated with each WriterSync.
	defaultBufferSize = 256 * 1024

	// defaultFlushInterval means the default flush interval
	defaultFlushInterval = 5 * time.Second
)

// defaultLogger is default logger for biz
// stdLogger is logger for std
// jupiterLogger is logger for jupiter framework
var defaultLogger, stdLogger, jupiterLogger *Logger

func init() {
	SetDefault(Config{
		Name:  "default",
		Debug: true,
	}.Build())

	SetJupiter(Config{
		Name:  "jupiter",
		Debug: true,
	}.Build())
}

func newLogger(config *Config) *zap.Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}

	zapOptions = append(zapOptions, zap.Hooks(hook))

	var ws zapcore.WriteSyncer
	if config.Debug || xdebug.IsDevelopmentMode() {
		ws = os.Stdout
	} else {
		ws = zapcore.AddSync(newRotate(config))
	}

	if config.Async {
		ws = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ws),
			FlushInterval: defaultFlushInterval,
			Size:          defaultBufferSize,
		}
		hooks.Register(hooks.Stage_AfterStop, func() { _ = ws.Sync() })
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	// encoderConfig := defaultZapConfig()
	// if config.Debug {
	// 	encoderConfig = defaultDebugConfig()
	// }
	encoderConfig := *config.EncoderConfig
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.Debug || xdebug.IsDevelopmentMode() {
					return zapcore.NewConsoleEncoder(encoderConfig)
				}
				return zapcore.NewJSONEncoder(encoderConfig)
			}(),
			ws,
			lv,
		)
	}

	zapLogger := zap.New(
		core,
		zapOptions...,
	)

	return zapLogger.Named(config.Name)
}

// DefaultZapConfig ...
func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// DebugEncodeLevel ...
func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = color.RedString
	switch lv {
	case zapcore.DebugLevel:
		colorize = color.BlueString
	case zapcore.InfoLevel:
		colorize = color.GreenString
	case zapcore.WarnLevel:
		colorize = color.YellowString
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = color.RedString
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}
