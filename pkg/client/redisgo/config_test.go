package redisgo

import (
	"bytes"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/util/xstring"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"

	"github.com/stretchr/testify/assert"
)

func TestStdConfig(t *testing.T) {
	var configStr = `
[jupiter.redisgo]
    [jupiter.redisgo.test.stub]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
			username=""
			password="123"

    [jupiter.redisgo.test.cluster]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config on addr nil", func(t *testing.T) {
		var config *Config
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "please set redisgo stub addr:test")
				assert.Nil(t, config)
			}
		}()
		config = StdStubConfig("test")
		assert.Nil(t, config) //不会执行到这里
	})
	t.Run("cluster config on addrs nil", func(t *testing.T) {
		var config *Config
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "please set cluster redisgo addrs:test")
				assert.Nil(t, config)
			}
		}()
		config = StdClusterConfig("test")
		assert.Nil(t, config) //不会执行到这里
	})
}
func TestNewConfig(t *testing.T) {
	var configStr = `
[jupiter.redisgo]
    [jupiter.redisgo.test.stub]
			addr = "redis://:user111:password222@127.0.0.1:6379"
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
			username=""
			password="123"

    [jupiter.redisgo.test.cluster]
			addrs = ["redis://:user111:password222@127.0.0.2:6379","redis://:user111:password222@127.0.0.1:6379"]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config", func(t *testing.T) {

		config := StdStubConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)
		assert.Equal(t, config.MaxRetries, 0)
		assert.Equal(t, config.EnableMetricInterceptor, true)
		assert.Equal(t, config.EnableTraceInterceptor, true)
		assert.Equal(t, config.EnableAccessLogInterceptor, false)
		assert.Equal(t, config.Debug, false)
		assert.Equal(t, config.Username, "")
		assert.Equal(t, config.Password, "123")
		assert.Equal(t, config.Addr, "127.0.0.1:6379")

	})
	t.Run("std cluster config", func(t *testing.T) {

		config := StdClusterConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)
		assert.Equal(t, config.MaxRetries, 0)
		assert.Equal(t, config.EnableMetricInterceptor, true)
		assert.Equal(t, config.EnableTraceInterceptor, true)
		assert.Equal(t, config.EnableAccessLogInterceptor, false)
		assert.Equal(t, config.Debug, false)
		assert.Equal(t, config.Username, "")
		assert.Equal(t, config.Password, "user111:password222")
		assert.Equal(t, config.Addr, "")
		assert.Equal(t, xstring.Json(config.Addrs), `["127.0.0.2:6379","127.0.0.1:6379"]`)

	})
}

func TestOldConfig(t *testing.T) {
	var configStr = `
[jupiter.redisgo]
    [jupiter.redisgo.test]
        [jupiter.redisgo.test.stub]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
            [jupiter.redisgo.test.stub.master]
                addr="redis://:user111:password222@127.0.0.1:6379"
            [jupiter.redisgo.test.stub.slaves]
                addr=[
                    "redis://:user111:password222@127.0.0.2:6379",
                ]
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config", func(t *testing.T) {

		config := StdStubConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)
		assert.Equal(t, config.MaxRetries, 0)
		assert.Equal(t, config.EnableMetricInterceptor, true)
		assert.Equal(t, config.EnableTraceInterceptor, true)
		assert.Equal(t, config.EnableAccessLogInterceptor, false)
		assert.Equal(t, config.Debug, false)
		assert.Equal(t, config.Username, "")
		assert.Equal(t, config.Password, "user111:password222")
		assert.Equal(t, config.Addr, "127.0.0.1:6379")

	})
	t.Run("std cluster config", func(t *testing.T) {

		config := StdClusterConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)
		assert.Equal(t, config.MaxRetries, 0)
		assert.Equal(t, config.EnableMetricInterceptor, true)
		assert.Equal(t, config.EnableTraceInterceptor, true)
		assert.Equal(t, config.EnableAccessLogInterceptor, false)
		assert.Equal(t, config.Debug, false)
		assert.Equal(t, config.Username, "")
		assert.Equal(t, config.Password, "user111:password222")
		assert.Equal(t, config.Addr, "")
		assert.Equal(t, xstring.Json(config.Addrs), `["127.0.0.2:6379","127.0.0.1:6379"]`)

	})
}
