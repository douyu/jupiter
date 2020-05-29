// Copyright 2020 Douyu
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

package xcron

import (
	"fmt"
	"runtime"
	"time"

	"github.com/douyu/jupiter/pkg/metric"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/robfig/cron/v3"
)

// StdConfig ...
func StdConfig(name string) Config {
	return RawConfig("jupiter.cron." + name)
}

// RawConfig ...
func RawConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal", xlog.String("key", key))
	}

	return config
}

// DefaultConfig ...
func DefaultConfig() Config {
	return Config{
		logger:          xlog.DefaultLogger,
		wrappers:        []JobWrapper{},
		WithSeconds:     false,
		ImmediatelyRun:  false,
		ConcurrentDelay: -1, // skip
	}
}

// Config ...
type Config struct {
	WithSeconds     bool
	ConcurrentDelay time.Duration
	ImmediatelyRun  bool

	wrappers []JobWrapper
	logger   *xlog.Logger
	parser   cron.Parser
}

// WithChain ...
func (config *Config) WithChain(wrappers ...JobWrapper) Config {
	if config.wrappers == nil {
		config.wrappers = []JobWrapper{}
	}
	config.wrappers = append(config.wrappers, wrappers...)
	return *config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) Config {
	config.logger = logger
	return *config
}

// WithParser ...
func (config *Config) WithParser(parser Parser) Config {
	config.parser = parser
	return *config
}

// Build ...
func (config Config) Build() *Cron {
	if config.WithSeconds {
		config.parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	if config.ConcurrentDelay > 0 { // 延迟
		config.wrappers = append(config.wrappers, delayIfStillRunning(config.logger))
	} else if config.ConcurrentDelay < 0 { // 跳过
		config.wrappers = append(config.wrappers, skipIfStillRunning(config.logger))
	} else {
		// 默认不延迟也不跳过
	}
	return newCron(&config)
}

type wrappedLogger struct {
	*xlog.Logger
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.Errorw("cron "+msg, append(keysAndValues, "err", err)...)
}

type wrappedJob struct {
	NamedJob
	logger *xlog.Logger
}

// Run ...
func (wj wrappedJob) Run() {
	tracer := xlog.NewTracer()
	metric.WorkerMetricsHandler.GetHandlerCounter().
		WithLabelValues("cron", wj.Name(), "begin").Inc()
	var beg = time.Now()
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			metric.WorkerMetricsHandler.GetHandlerCounter().
				WithLabelValues("cron", wj.Name(), "over err").Inc()
			tracer.Error(
				xlog.Any("err", err),
				xlog.String("event", "recover"),
				xlog.String("stack", string(buf)),
			)
		} else {
			metric.WorkerMetricsHandler.GetHandlerCounter().
				WithLabelValues("cron", wj.Name(), "over suc").Inc()
		}
		metric.WorkerMetricsHandler.GetHandlerHistogram().
			WithLabelValues("cron", wj.Name()).Observe(time.Since(beg).Seconds())
		tracer.Info(
			xlog.String("name", wj.Name()),
		)
		tracer.Flush("run job", wj.logger)
	}()
	if err := wj.NamedJob.Run(); err != nil {

	}
}
