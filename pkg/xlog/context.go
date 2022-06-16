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
	loggerKey  struct{}
	traceIDKey struct{}
)

func NewContext(ctx context.Context, l *Logger) context.Context {
	traceID := GetTraceID(ctx)
	if traceID == "" {
		return context.WithValue(ctx, loggerKey{}, l)
	}
	return context.WithValue(ctx, loggerKey{}, l.With(String(traceIDField, traceID)))
}

func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerKey{}).(*Logger)
	if !ok {
		return defaultLogger // default logger
	}
	return l
}

func SetTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return ""
	}
	return traceID
}
