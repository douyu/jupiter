package redis

import (
	"bytes"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"

	"github.com/douyu/jupiter/pkg/core/constant"
)

func TestStdConfig(t *testing.T) {
	assert.Equal(t, constant.ConfigKey("redis", "test", "stub"), "jupiter.redis.test.stub")

	var configStr = `
[jupiter.redis]
    [jupiter.redis.test.stub]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
            username=""
            password="123"

	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config on addr nil", func(t *testing.T) {
		var config *Config
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "no master or slaves addr set:jupiter.redis.test.stub")
				assert.Nil(t, config)
			}
		}()
		config = StdConfig("test")
		assert.Nil(t, config) //不会执行到这里
	})

}

func TestConfig(t *testing.T) {
	var configStr = `
[jupiter.redis]
    [jupiter.redis.test]
        [jupiter.redis.test.stub]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
            [jupiter.redis.test.stub.master]
                addr="redis://:user111:password222@127.0.0.1:6379"
            [jupiter.redis.test.stub.slaves]
                addr=[
                    "redis://:user111:password222@127.0.0.2:6379",
                ]
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("std config", func(t *testing.T) {

		config := StdConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)
		assert.Equal(t, config.MaxRetries, 0)
		assert.Equal(t, config.EnableMetricInterceptor, true)
		assert.Equal(t, config.EnableTraceInterceptor, true)
		assert.Equal(t, config.EnableAccessLogInterceptor, false)
		assert.Equal(t, config.Debug, false)

		assert.Equal(t, config.Master.Addr, "redis://:user111:password222@127.0.0.1:6379")
		assert.Equal(t, len(config.Slaves.Addr), 2)

	})

}
