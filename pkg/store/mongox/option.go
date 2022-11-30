package mongox

import (
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Config ...
type (
	Config struct {
		Name string
		// DSN地址
		DSN string `json:"dsn" toml:"dsn"`
		// 创建连接的超时时间
		SocketTimeout time.Duration `json:"socketTimeout" toml:"socketTimeout"`
		// 连接池大小(最大连接数)
		PoolLimit int `json:"poolLimit" toml:"poolLimit"`

		logger *zap.Logger
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("mongo." + name))
}

// RawConfig 裸配置
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config.Name = key
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		panic(err)
	}
	return config
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		SocketTimeout: cast.ToDuration("5s"),
		PoolLimit:     100,
		logger:        xlog.Jupiter().Named(ecode.ModStoreMongo),
	}
}

func (config *Config) Build() *mongo.Client {
	return newSession(config)
}

func (config *Config) Singleton() *mongo.Client {
	if val, ok := singleton.Load(constant.ModuleStoreMongoDB, config.Name); ok {
		return val.(*mongo.Client)
	}

	val := config.Build()
	singleton.Store(constant.ModuleStoreMongoDB, config.Name, val)

	return val
}
