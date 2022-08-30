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
	"go.uber.org/zap"
)

// Jupiter returns framework logger
func Jupiter() *zap.Logger {
	return jupiterLogger
}

// Default returns default logger
func Default() *zap.Logger {
	return defaultLogger
}

func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	defaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

func With(fields ...Field) *Logger {
	return defaultLogger.With(fields...)
}

func WithOptions(opts ...Option) *Logger {
	return defaultLogger.WithOptions(opts...)
}

func Named(s string) *Logger {
	return defaultLogger.Named(s)
}
