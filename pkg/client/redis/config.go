package redis

import (
	"strings"
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config ...
type Config struct {
	// Master host:port addresses of Master node
	Master struct {
		Addr string `json:"addr" toml:"addr"`
	} `json:"master" toml:"master"`
	// Slaves A list of host:port addresses of Slave nodes.
	Slaves struct {
		Addr []string `json:"addr" toml:"addr"`
	} `json:"slaves" toml:"slaves"`

	Addr     []string `json:"addr" toml:"addr"`
	Username string   `json:"username" toml:"username"`
	Password string   `json:"password" toml:"password"`

	/****** for github.com/go-redis/redis/v8 ******/
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
	ReadOnMaster bool `json:"readOnMaster" toml:"readOnMaster"`
	// nice option
	Debug bool `json:"debug" toml:"debug"`
	// a require will be recorded if cost bigger than this
	SlowLogThreshold time.Duration `json:"slowThreshold" toml:"slowThreshold"`
	// EnableMetric .. default true
	EnableMetricInterceptor bool `json:"enableMetric" toml:"enableMetric"`
	// EnableTrace .. default true
	EnableTraceInterceptor bool `json:"enableTrace" toml:"enableTrace"`
	// EnableAccessLog .. default false
	EnableAccessLogInterceptor bool `json:"enableAccessLog" toml:"enableAccessLog"`
	// EnableSentinel .. default true
	EnableSentinel bool `json:"enableSentinel" toml:"enableSentinel"`
	// OnDialError panic|error
	OnDialError string `json:"level"`
	logger      *zap.Logger
	name        string
}

// DefaultConfig default config ...
func DefaultConfig() *Config {
	return &Config{
		name:                    "default",
		DB:                      0,
		PoolSize:                200,
		MinIdleConns:            20,
		DialTimeout:             cast.ToDuration("3s"),
		ReadTimeout:             cast.ToDuration("1s"),
		WriteTimeout:            cast.ToDuration("1s"),
		IdleTimeout:             cast.ToDuration("60s"),
		ReadOnMaster:            true,
		Debug:                   false,
		EnableMetricInterceptor: true,
		EnableTraceInterceptor:  true,
		EnableSentinel:          true,
		SlowLogThreshold:        cast.ToDuration("250ms"),
		logger:                  xlog.Jupiter().Named(ecode.ModClientRedis),
		OnDialError:             "panic",
	}
}

// StdConfig 哨兵模式，支持主从结构redis集群
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("redis", name, "stub"))
}

func RawConfig(key string) *Config {
	var config = DefaultConfig()

	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		config.logger.Panic("unmarshal config:"+key, xlog.FieldErr(err), xlog.FieldName(key), xlog.FieldExtMessage(config))
	}

	if config.Master.Addr != "" && config.ReadOnMaster {
		config.Slaves.Addr = append(config.Slaves.Addr, config.Master.Addr)
	}
	if config.Master.Addr == "" && len(config.Slaves.Addr) == 0 {
		config.logger.Panic("no master or slaves addr set:"+key, xlog.FieldName(key), xlog.FieldExtMessage(config))
	}
	config.name = key
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}

	return config

}
func getUsernameAndPassword(addr string) (realAddr string, username, password string) {
	addr = strings.TrimPrefix(addr, "redis://")
	addr = strings.TrimPrefix(addr, "rediss://")
	arr := strings.Split(addr, "@")
	if len(arr) < 2 {
		return addr, "", ""
	}
	realAddr = arr[1]

	// username:password
	auth := arr[0]
	subArr := strings.Split(auth, ":")
	if len(subArr) >= 2 {
		username = subArr[0]
		password = strings.Join(subArr[1:], ":")
	}
	return realAddr, username, password
}
