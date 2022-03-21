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
		return DefaultLogger // default logger
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
