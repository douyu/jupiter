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

package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

type mqLogger struct {
	logger *zap.Logger
}

const (
	defaultCallerSkip = 2
)

func init() {
	hooks.Register(hooks.Stage_AfterLoadConfig, func() {
		SetLogger(xlog.Jupiter())
	})
}

func SetLogger(logger *zap.Logger) {
	rlog.SetLogLevel("debug")
	rlog.SetLogger(&mqLogger{
		logger: logger.Named(ecode.ModeClientRocketMQ).
			WithOptions(zap.AddCallerSkip(defaultCallerSkip)),
	})
}

func (l *mqLogger) Debug(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	fs := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		fs = append(fs, zap.Any(key, value))
	}

	l.logger.Debug(msg, fs...)
}

func (l *mqLogger) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	fs := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		fs = append(fs, zap.Any(key, value))
	}

	// Here we reguard the info level as debug level
	l.logger.Debug(msg, fs...)
}

func (l *mqLogger) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	fs := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		fs = append(fs, zap.Any(key, value))
	}

	l.logger.Warn(msg, fs...)
}

func (l *mqLogger) Error(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	fs := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		fs = append(fs, zap.Any(key, value))
	}

	l.logger.Error(msg, fs...)
}

func (l *mqLogger) Fatal(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	fs := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		fs = append(fs, zap.Any(key, value))
	}

	l.logger.Fatal(msg, fs...)
}

func (l *mqLogger) Level(level string) {

}

func (l *mqLogger) OutputPath(path string) (err error) {
	return nil
}
