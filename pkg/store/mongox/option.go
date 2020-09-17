package mongox

import (
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"go.mongodb.org/mongo-driver/mongo"
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
	}
)

// StdConfig 返回标准配置
func StdConfig(name string) Config {
	return RawConfig("jupiter.mongo." + name)
}

// RawConfig 裸配置
// example: minerva.mongodb.demo
func RawConfig(key string) Config {
	var config = DefaultConfig()
	config.Name = key
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		panic(err)
	}
	return config
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		SocketTimeout: xtime.Duration("5s"),
		PoolLimit:     100,
	}
}

func (config Config) Build() *mongo.Client {
	return newSession(config)
}
