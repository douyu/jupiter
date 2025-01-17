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

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/core/sentinel"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var errSlowCommand = errors.New("http resty slow command")

// Client ...
type Client = resty.Client

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
		// 配置名称
		Name string `json:"name"`
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("resty." + name))
}

// RawConfig 返回配置
func RawConfig(key string) *Config {
	config := DefaultConfig()
	config.Name = key

	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("unmarshal config", xlog.FieldName(key), xlog.FieldExtMessage(config))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}
	return config
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Debug:            false,
		EnableMetric:     true,
		EnableTrace:      true,
		RetryCount:       0,
		RetryWaitTime:    cast.ToDuration("100ms"),
		RetryMaxWaitTime: cast.ToDuration("100ms"),
		Addr:             "",
		SlowThreshold:    cast.ToDuration("500ms"),
		Timeout:          cast.ToDuration("3000ms"),
		EnableAccessLog:  false,
		EnableSentinel:   true,
		logger:           xlog.Jupiter().Named(ecode.ModeClientResty),
		Name:             "default",
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

		if config.EnableSentinel {
			entry := sentinel.FromContext(r.Context())
			if entry != nil {
				entry.Exit(base.WithError(err))
			}
		}
	})

	tracer := xtrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("http"),
	}

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		if config.EnableTrace {

			ctx, _ := tracer.Start(r.Context(), r.URL, propagation.HeaderCarrier(r.Header), trace.WithAttributes(attrs...))

			r.SetContext(ctx)
		}

		if config.EnableSentinel {
			entry, err := sentinel.Entry(r.URL, api.WithTrafficType(base.Outbound), api.WithResourceType(base.ResTypeWeb))
			if err != nil {
				return err
			}

			r.SetContext(sentinel.WithContext(r.Context(), entry))
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
			span.SetAttributes(semconv.HTTPClientAttributesFromHTTPRequest(r.Request.RawRequest)...)
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int64(int64(r.StatusCode())),
			)

			if r.IsError() {
				span.RecordError(errors.New(r.Status()))
				span.SetStatus(codes.Error, r.Status())
			}

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

		if config.EnableSentinel {
			entry := sentinel.FromContext(r.Request.Context())
			if entry != nil {
				entry.Exit()
			}
		}

		return nil
	})

	return client, nil
}

// Singleton returns a singleton client conn.
func (config *Config) Singleton() (*Client, error) {
	if client, ok := singleton.Load(constant.ModuleClientResty, config.Name); ok && client != nil {
		return client.(*Client), nil
	}

	client, err := config.Build()
	if err != nil {
		xlog.Jupiter().Error("build resty client failed", zap.Error(err))
		return nil, err
	}

	singleton.Store(constant.ModuleClientResty, config.Name, client)

	return client, nil
}

// MustBuild panics when error found.
func (c *Config) MustBuild() *resty.Client {
	return lo.Must(c.Build())
}

// MustSingleton panics when error found.
func (config *Config) MustSingleton() *Client {
	return lo.Must(config.Singleton())
}
