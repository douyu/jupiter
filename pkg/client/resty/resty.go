package resty

import (
	"errors"
	"net/http"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

var _logger = xlog.DefaultLogger.With(xlog.FieldMod("resty"))
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
		// 影子流量开关
		ShadowSwitch string // on打开， off关闭， watch观察者模式（关闭且打印影子日志）
		// 默认noHostUrl to mockRes, eg: /test/v1=>"{\"rid\":20}"
		MockResMap map[string]struct {
			Url  string `json:"url" toml:"url"`
			Data string `json:"data" toml:"data"`
		} `json:"mockRes" toml:"mockRes"`
		// 所有方法默认返回的mock数据
		DefaultMockRes string                   `json:"defaultMockRes" toml:"defaultMockRes"`
		RetryCondition resty.RetryConditionFunc `json:"-" toml:"-"`
		// 隐藏 X-DY header
		DisableDYHeader bool `json:"disableDYHeader" toml:"disableDYHeader"`
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) Config {
	return RawConfig("minerva.resty." + name)
}

// RawConfig 返回配置
// example: RawConfig("minerva.http.demo")
func RawConfig(key string) Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		xlog.Panic("unmarshal config", xlog.FieldName(key), xlog.FieldExtMessage(config))
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
		Timeout:          xtime.Duration("1000ms"),
		EnableAccessLog:  false,
		EnableSentinel:   true,
		ShadowSwitch:     "off",
		DefaultMockRes:   "",
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
			span := opentracing.SpanFromContext(r.Context())
			ext.LogError(span, err)
			span.Finish()
		}

		if config.EnableMetric {
			metric.ClientHandleCounter.WithLabelValues(metric.TypeHTTP, "resty", r.Method, r.RawRequest.Host, "error").Inc()
		}
	})

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		if config.EnableTrace {
			// 设置trace
			span, ctx := trace.StartSpanFromContext(
				r.Context(),
				"http_client_request",
				trace.MetadataExtractor(r.Header),
			)

			ext.SpanKindRPCClient.Set(span)
			ext.Component.Set(span, "http")
			ext.HTTPUrl.Set(span, r.URL)
			ext.HTTPMethod.Set(span, r.Method)

			r.SetContext(trace.HeaderInjector(ctx, r.Header))
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
			trace.SpanFromContext(r.Request.Context()).Finish()
		}

		if config.SlowThreshold > time.Duration(0) {
			// 慢日志
			if cost > config.SlowThreshold {
				_logger.Error("slow",
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
		xlog.Panic("resty build failed", zap.Error(err))
	}

	return cc
}
