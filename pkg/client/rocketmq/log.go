package rocketmq

import (
	"go.uber.org/zap"
)

type mqLogger struct {
	logger *zap.Logger
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

	l.logger.Info(msg, fs...)
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

	fs := []zap.Field{}

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
