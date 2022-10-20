package xlog

import (
	"github.com/douyu/jupiter/pkg/core/metric"
	"go.uber.org/zap/zapcore"
)

// hook does capture metrics like log number group by level, etc ...
func hook(e zapcore.Entry) error {
	metric.LogLevelCounter.WithLabelValues(e.LoggerName, e.Level.String()).Inc()

	return nil
}
