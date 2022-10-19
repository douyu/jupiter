package redisgo

import (
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"

	"go.uber.org/zap"

	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xdebug"

	"git.dz11.com/vega/minerva/util/xtime"
)

// StubConfig are used to configure a stub client and should be
// passed to NewStubClient.
type Config struct {
	// Master host:port addresses of Master node
	Addr string `json:"addr" toml:"addr"`
	// Slaves A list of host:port addresses of Slave nodes.
	Addrs []string `json:"addrs" toml:"addrs"`

	/****** for github.com/go-redis/redis/v8 ******/
	Username string `json:"username" toml:"username"`
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
	// Enables read-only commands on slave nodes.
	ReadOnly bool

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
	EnableAccessLogInterceptor bool `json:"enableAccessLog" toml:"enableAccessLog"`
	// OnDialError panic|error
	OnDialError string `json:"level"`
	logger      *zap.Logger
	name        string
}

// AddrString 获取地址字符串, 用于 log, metric, trace 中的 label
func (config *Config) AddrString() string {
	addr := config.Addr
	if len(config.Addrs) > 0 {
		addr = strings.Join(config.Addrs, ",")
	}
	return addr
}

// DefaultConfig default config ...
func DefaultConfig() *Config {
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
		OnDialError:             "panic",
	}
}

// StdStubConfig ...
func StdStubConfig(name string) *Config {
	var config = DefaultConfig()
	key := "jupiter.redisgo." + name + ".stub"
	if !strings.HasSuffix(key, ".stub") {
		key = key + ".stub"
	}

	if err := cfg.UnmarshalKey(key, &config, cfg.TagName("toml")); err != nil {
		config.logger.Panic("unmarshal config:"+key, xlog.FieldErr(err), xlog.FieldName(key), xlog.FieldExtMessage(config))
	}
	if config.Addr == "" { // 兼容master的写法
		config.Addr = cfg.GetString(key + ".master.addr")
	}
	addr, user, pass := getUsernameAndPassword(config.Addr)
	if user != "" && config.Username == "" { // 配置里的username优先
		config.Username = user
	}
	if pass != "" && config.Password == "" { // 配置里的password优先
		config.Password = pass
	}
	if addr != "" {
		config.Addr = addr
	}
	if config.Addr == "" {
		config.logger.Panic("please set redisgo stub addr:"+name, xlog.FieldName(key), xlog.FieldExtMessage(config))

	}
	config.name = name
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}

	return config

}

// StdClusterConfig ...
func StdClusterConfig(name string) *Config {
	var config = DefaultConfig()
	key := "jupiter.redisgo." + name

	if err := cfg.UnmarshalKey(key+".cluster", &config, cfg.TagName("toml")); err != nil {
		config.logger.Info("unmarshal cluster config:"+key+".cluster", xlog.FieldErr(err), xlog.FieldName(key), xlog.FieldExtMessage(config))

		if err2 := cfg.UnmarshalKey(key+".stub", &config, cfg.TagName("toml")); err2 != nil {
			config.logger.Panic("unmarshal cluster config:"+key+".stub", xlog.FieldErr(err2), xlog.FieldName(key), xlog.FieldExtMessage(config))
		}
		if len(config.Addrs) == 0 { // 兼容slaves的写法
			if oldAddr := cfg.GetStringSlice(key + ".stub.slaves.addr"); len(oldAddr) > 0 {
				config.Addrs = oldAddr
			}
		}
		if addr := cfg.GetString(key + ".stub.master.addr"); addr != "" {
			config.Addrs = append(config.Addrs, addr)
		}

	}
	if len(config.Addrs) <= 0 {
		config.logger.Panic("please set cluster redisgo addrs:"+config.name, xlog.FieldName(key), xlog.FieldExtMessage(config))
	}
	// 解析第一个地址
	addrs := []string{}
	for _, item := range config.Addrs {
		addr, user, pass := getUsernameAndPassword(item)
		if user != "" && config.Username == "" { // 配置里的username优先
			config.Username = user
		}
		if pass != "" && config.Password == "" { // 配置里的password优先
			config.Password = pass
		}
		addrs = append(addrs, addr)
	}
	config.Addrs = addrs

	config.name = name
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
