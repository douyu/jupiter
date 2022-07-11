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

package resty

import (
	"errors"
	"net/http"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/douyu/jupiter/pkg/xtrace"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var errSlowCommand = errors.New("http resty slow command")

// Config ...
type (
	// Config options
	Config struct {
		// Debug 开关
		Debug bool `json:"debug" toml:"debug"` // debug开关
		// 指标采集开关
		EnableMetric bool `json:"enableMetric" toml:"enableMetric"` // 指标采集开关
		// 链路追踪开关
		EnableTrace bool `json:"enableTrace" toml:"enableTrace"` // 链路开关
		// 失败重试次数
		RetryCount int `json:"retryCount" toml:"retryCount"` // 重试次数
		// 失败重试的间隔时间
		RetryWaitTime time.Duration `json:"retryWaitTime" toml:"retryWaitTime"` // 重试间隔时间
		// 失败重试的最贱等待时间
		RetryMaxWaitTime time.Duration `json:"retryMaxWaitTime" toml:"retryMaxWaitTime"` // 重试最大间隔时间
		// 目标服务地址
		Addr string `json:"addr" toml:"addr"` // 目标地址
		// 请求超时时间
		Timeout time.Duration `json:"timeout" toml:"timeout" `
		// 收到响应以后是否立即关闭连接
		CloseConnection bool `json:"closeConnection" toml:"closeConnection" `
		// 慢日志阈值
		SlowThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"` // slowlog 时间阈
		// 访问日志开关
		EnableAccessLog bool `json:"enableAccessLog" toml:"enableAccessLog"`
		// 熔断降级
		EnableSentinel bool `json:"enableSentinel" toml:"enableSentinel"`
		// 重试
		RetryCondition resty.RetryConditionFunc `json:"-" toml:"-"`
		// 日志
		logger *zap.Logger
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) Config {
	return RawConfig("jupiter.resty." + name)
}

// RawConfig 返回配置
func RawConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("unmarshal config", xlog.FieldName(key), xlog.FieldExtMessage(config))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}
	return config
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Debug:            false,
		EnableMetric:     true,
		EnableTrace:      true,
		RetryCount:       0,
		RetryWaitTime:    xtime.Duration("100ms"),
		RetryMaxWaitTime: xtime.Duration("100ms"),
		Addr:             "",
		SlowThreshold:    xtime.Duration("500ms"),
		Timeout:          xtime.Duration("3000ms"),
		EnableAccessLog:  false,
		EnableSentinel:   true,
		logger:           xlog.Jupiter().With(xlog.FieldMod("resty")),
	}
}

func (config *Config) Build() (*resty.Client, error) {
	if config.Addr == "" {
		return nil, errors.New("no addr found")
	}

	client := resty.New()
	client.SetBaseURL(config.Addr)
	client.SetTimeout(config.Timeout)
	client.SetDebug(config.Debug)
	client.SetRetryCount(config.RetryCount)
	client.SetCloseConnection(config.CloseConnection)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	client.AddRetryCondition(config.RetryCondition)

	if config.RetryWaitTime != time.Duration(0) {
		client.SetRetryWaitTime(config.RetryWaitTime)
	}

	if config.RetryMaxWaitTime != time.Duration(0) {
		client.SetRetryMaxWaitTime(config.RetryMaxWaitTime)
	}

	client.OnError(func(r *resty.Request, err error) {
		if config.EnableTrace {
			span := trace.SpanFromContext(r.Context())
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}
		if config.EnableMetric {
			metric.ClientHandleCounter.WithLabelValues(metric.TypeHTTP, "resty", r.Method, r.RawRequest.Host, "error").Inc()
		}

	})

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		if config.EnableTrace {
			tracer := xtrace.NewTracer(trace.SpanKindClient)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("http"),
			}
			ctx, span := tracer.Start(r.Context(), r.Method, propagation.HeaderCarrier(r.Header), trace.WithAttributes(attrs...))
			span.SetAttributes(
				semconv.RPCSystemKey.String("http"),
				semconv.PeerServiceKey.String("http_client_request"),
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPURLKey.String(r.URL),
			)
			r.SetContext(ctx)
		}
		return nil
	})

	client.SetPreRequestHook(func(c *resty.Client, r *http.Request) error {

		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {

		cost := r.Time()

		if config.EnableMetric {
			metric.ClientHandleCounter.WithLabelValues(metric.TypeHTTP, "resty", r.Request.Method, c.HostURL, r.Status()).Inc()
			metric.ClientHandleHistogram.WithLabelValues(metric.TypeHTTP, "resty", r.Request.Method, c.HostURL).Observe(cost.Seconds())
		}

		if config.EnableTrace {
			span := trace.SpanFromContext(r.Request.Context())
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int64(int64(r.StatusCode())),
			)
			span.End()
		}

		if config.SlowThreshold > time.Duration(0) {
			// 慢日志
			if cost > config.SlowThreshold {
				config.logger.Error("slow",
					xlog.FieldErr(errSlowCommand),
					xlog.FieldMethod(r.Request.Method),
					xlog.FieldCost(cost),
					xlog.FieldAddr(r.Request.URL),
					xlog.FieldCode(int32(r.StatusCode())),
				)
			}
		}

		return nil
	})

	return client, nil
}

func (c *Config) MustBuild() *resty.Client {
	cc, err := c.Build()
	if err != nil {
		xlog.Jupiter().Panic("resty build failed", zap.Error(err), zap.Any("config", c))
	}

	return cc
}
