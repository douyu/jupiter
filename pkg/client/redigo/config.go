package redigo

import (
	"context"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis"

	"go.uber.org/zap"

	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xdebug"

	"git.dz11.com/vega/minerva/util/xtime"
)

// StubConfig are used to configure a stub client and should be
// passed to NewStubClient.
type Config struct {
	// Master host:port addresses of Master node
	Addr string `json:"master" toml:"master"`
	// Slaves A list of host:port addresses of Slave nodes.
	Addrs []string `json:"slaves" toml:"slaves"`

	/****** for github.com/go-redis/redis/v8 ******/
	Password string `json:"password" toml:"password"`
	// DB default 0,not recommend
	DB int `json:"db" toml:"db"`
	// PoolSize applies per Stub node and not for the whole Stub.
	PoolSize int `json:"poolSize" toml:"poolSize"`
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `json:"maxRetries" toml:"maxRetries"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `json:"minIdleConns" toml:"minIdleConns"`
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value 0 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration `json:"readTimeout" toml:"readTimeout"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration `json:"writeTimeout" toml:"writeTimeout"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration `json:"idleTimeout" toml:"idleTimeout"`

	/****** for jupiter ******/
	// nice option
	Debug bool `json:"debug" toml:"debug"`
	// a require will be recorded if cost bigger than this
	SlowLogThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// EnableMetric .. default true
	EnableMetricInterceptor bool `json:"enableMetric" toml:"enableMetric"`
	// EnableTrace .. default true
	EnableTraceInterceptor bool `json:"enableTrace" toml:"enableTrace"`
	// EnableAccessLog .. default false
	EnableAccessLogInterceptor bool        `json:"enableAccessLog" toml:"enableAccessLog"`
	logger                     *zap.Logger //
}

// AddrString 获取地址字符串, 用于 log, metric, trace 中的 label
func (c *Config) AddrString() string {
	addr := c.Addr
	if len(c.Addrs) > 0 {
		addr = strings.Join(c.Addrs, ",")
	}
	return addr
}

// DefaultStubConfig default config ...
func DefaultStubConfig() *Config {
	return &Config{
		DB:                      0,
		PoolSize:                200,
		MinIdleConns:            20,
		DialTimeout:             xtime.Duration("3s"),
		ReadTimeout:             xtime.Duration("1s"),
		WriteTimeout:            xtime.Duration("1s"),
		IdleTimeout:             xtime.Duration("60s"),
		Debug:                   false,
		EnableMetricInterceptor: true,
		EnableTraceInterceptor:  true,
		SlowLogThreshold:        xtime.Duration("250ms"),
		logger:                  xlog.Jupiter().With(xlog.FieldMod("redigo")),
	}
}

// StdStubConfig ...
func StdStubConfig(name string) (string, *Config) {
	return name, RawStubConfig("jupiter.redisgo." + name + ".stub")
}

// RawStubConfig ...
func RawStubConfig(key string) *Config {
	var config = DefaultStubConfig()
	if !strings.HasSuffix(key, ".stub") {
		key = key + ".stub"
	}

	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		config.logger.Panic("unmarshal config", xlog.FieldErr(err), xlog.FieldName(key), xlog.FieldExtMessage(config))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}

	return config
}
func (config *Config) Build() {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})

	for _, incpt := range config.interceptors {
		stubClient.AddHook(incpt)
	}

	if err := stubClient.Ping(context.Background()).Err(); err != nil {
		switch c.config.OnFail {
		case "panic":
			c.logger.Panic("start stub redis", elog.FieldErr(err))
		default:
			c.logger.Error("start stub redis", elog.FieldErr(err))
		}
	}
	return stubClient
}
