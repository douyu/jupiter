package xgrpclog

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
)

const (
	defaultCallerSkip = 4

	// See https://github.com/grpc/grpc-go/blob/v1.35.0/grpclog/loggerv2.go#L77-L86
	grpcLvlInfo  = 0
	grpcLvlWarn  = 1
	grpcLvlError = 2
	grpcLvlFatal = 3
)

var (
	grpcToZapLevel = map[int]zapcore.Level{
		grpcLvlInfo:  zapcore.DebugLevel,
		grpcLvlWarn:  zapcore.WarnLevel,
		grpcLvlError: zapcore.ErrorLevel,
		grpcLvlFatal: zapcore.FatalLevel,
	}
)

func init() {
	hooks.Register(hooks.Stage_AfterLoadConfig, func() {
		SetLogger(xlog.Jupiter())
	})
}

// SetLogger sets loggerWrapper to grpclog
func SetLogger(logger *xlog.Logger) {
	logger = logger.Named(ecode.ModClientGrpc).WithOptions(zap.AddCallerSkip(defaultCallerSkip))
	grpclog.SetLoggerV2(&loggerWrapper{logger: logger, sugar: logger.Sugar()})
}

// loggerWrapper wraps xlog.Logger into a LoggerV2.
type loggerWrapper struct {
	logger *xlog.Logger
	sugar  *zap.SugaredLogger
}

// Info logs to INFO log
func (l *loggerWrapper) Info(args ...interface{}) {
	l.logger.Info(sprint(args...))
}

// Infoln logs to INFO log
func (l *loggerWrapper) Infoln(args ...interface{}) {
	if l.logger.Core().Enabled(grpcToZapLevel[grpcLvlInfo]) {
		l.logger.Info(sprint(args...))
	}
}

// Infof logs to INFO log
func (l *loggerWrapper) Infof(format string, args ...interface{}) {
	l.sugar.Infof(sprintf(format, args...))
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warning(args ...interface{}) {
	l.logger.Warn(sprint(args...))
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warningln(args ...interface{}) {
	if l.logger.Core().Enabled(grpcToZapLevel[grpcLvlWarn]) {
		l.logger.Warn(sprint(args...))
	}
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warningf(format string, args ...interface{}) {
	l.logger.Warn(sprintf(format, args...))
}

// Error logs to ERROR log
func (l *loggerWrapper) Error(args ...interface{}) {
	l.logger.Error(sprint(args...))
}

// Errorn logs to ERROR log
func (l *loggerWrapper) Errorln(args ...interface{}) {
	if l.logger.Core().Enabled(grpcToZapLevel[grpcLvlError]) {
		l.logger.Error(sprint(args...))
	}
}

// Errorf logs to ERROR log
func (l *loggerWrapper) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(sprintf(format, args...))
}

// Fatal logs to ERROR log
func (l *loggerWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(sprint(args...))
}

// Fatalln logs to ERROR log
func (l *loggerWrapper) Fatalln(args ...interface{}) {
	l.logger.Fatal(sprint(args...))
}

// Error logs to ERROR log
func (l *loggerWrapper) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(sprintf(format, args...))
}

// v returns true for all verbose level.
func (l *loggerWrapper) V(v int) bool {
	return l.logger.Core().Enabled(grpcToZapLevel[v])
}

func sprint(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
