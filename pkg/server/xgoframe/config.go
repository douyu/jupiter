package xgoframe

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/pkg/errors"
)

// HTTP config
type Config struct {
	Host          string
	Port          int
	Debug         bool
	DisableMetric bool
	DisableTrace  bool

	SlowQueryThresholdInMilli int64

	logger *xlog.Logger
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                      "127.0.0.1",
		Port:                      8099,
		Debug:                     false,
		SlowQueryThresholdInMilli: 500, // 500ms
		logger:                    xlog.JupiterLogger.With(xlog.FieldMod("server.goframe")),
	}
}

// Jupiter Standard HTTP Server config
func StdConfig(name string) *Config {
	return RawConfig("jupiter.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil &&
		errors.Cause(err) != conf.ErrInvalidKey {
		config.logger.Panic("http server parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// WithHost ...
func (config *Config) WithHost(host string) *Config {
	config.Host = host
	return config
}

// WithPort ...
func (config *Config) WithPort(port int) *Config {
	config.Port = port
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() *Server {
	serve := newServer(config)

	serve.BindMiddlewareDefault()
	serve.Use(recoverMiddleware(config.logger, config.SlowQueryThresholdInMilli))
	//
	if !config.DisableMetric {
		serve.Use(metricServerInterceptor())
	}
	if !config.DisableTrace {
		serve.Use(traceServerInterceptor())
	}
	return serve
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
