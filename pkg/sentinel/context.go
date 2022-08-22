package sentinel

import (
	"context"

	"github.com/alibaba/sentinel-golang/core/base"
)

var (
	sentinelEntryKey = struct{}{}
)

func WithContext(ctx context.Context, val *base.SentinelEntry) context.Context {
	return context.WithValue(ctx, sentinelEntryKey, val)
}

func FromContext(ctx context.Context) *base.SentinelEntry {
	if ctx == nil {
		return nil
	}
	if val, ok := ctx.Value(sentinelEntryKey).(*base.SentinelEntry); ok {
		return val
	}
	return nil
}
