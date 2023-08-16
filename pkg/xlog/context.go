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
	"context"
)

const (
	traceIDField = "trace_id"
)

type (
	defaultLoggerKey struct{}
	jupiterLoggerKey struct{}
)

func NewContext(ctx context.Context, l *Logger, traceID string) context.Context {
	if l == jupiterLogger {
		return newContextWithJupiterLogger(ctx, l, traceID)
	}
	return newContextWithDefaultLogger(ctx, l, traceID)
}

func newContextWithDefaultLogger(ctx context.Context, l *Logger, traceID string) context.Context {
	return context.WithValue(ctx, defaultLoggerKey{}, l.With(String(traceIDField, traceID)))
}

func newContextWithJupiterLogger(ctx context.Context, l *Logger, traceID string) context.Context {
	return context.WithValue(ctx, jupiterLoggerKey{}, l.With(String(traceIDField, traceID)))
}

// Deprecated: use xlog.L instead
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return defaultLogger
	}

	l, ok := ctx.Value(defaultLoggerKey{}).(*Logger)
	if !ok {
		return defaultLogger // default logger
	}
	return l
}

func getDefaultLoggerFromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return defaultLogger
	}

	l, ok := ctx.Value(defaultLoggerKey{}).(*Logger)
	if !ok {
		return defaultLogger // default logger
	}
	return l
}

func getJupiterLoggerFromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return jupiterLogger
	}

	l, ok := ctx.Value(jupiterLoggerKey{}).(*Logger)
	if !ok {
		return jupiterLogger // jupiter logger
	}
	return l
}
