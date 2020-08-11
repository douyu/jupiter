package xgrpclog

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc/grpclog"
)

// SetLogger sets loggerWrapper to grpclog
func SetLogger(logger *xlog.Logger) {
	grpclog.SetLoggerV2(&loggerWrapper{logger: logger})
}

// loggerWrapper wraps xlog.Logger into a LoggerV2.
type loggerWrapper struct {
	logger *xlog.Logger
}

// Info logs to INFO log
func (l *loggerWrapper) Info(args ...interface{}) {
	l.logger.Info(sprint(args...))
}

// Infoln logs to INFO log
func (l *loggerWrapper) Infoln(args ...interface{}) {
	l.logger.Info(sprint(args...))
}

// Infof logs to INFO log
func (l *loggerWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(sprintf(format, args...))
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warning(args ...interface{}) {
	l.logger.Warn(sprint(args...))
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warningln(args ...interface{}) {
	l.logger.Warn(sprint(args...))
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
	l.logger.Error(sprint(args...))
}

// Errorf logs to ERROR log
func (l *loggerWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(sprintf(format, args...))
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
	l.logger.Fatalf(sprintf(format, args...))
}

// v returns true for all verbose level.
func (l *loggerWrapper) V(v int) bool {
	return true
}

func sprint(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
