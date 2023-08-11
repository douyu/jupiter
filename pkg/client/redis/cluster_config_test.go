package redis

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClusterOption_Nil(t *testing.T) {
	assert.Equal(t, constant.ConfigKey("redis", "test", "cluster"), "jupiter.redis.test.cluster")

	var configStr = `
[jupiter.redis.test]
    [jupiter.redis.test.cluster]
            dialTimeout="2s"
            readTimeout="5s"
            idleTimeout="60s"
            username="root"
            password="123"
			addr = ["r-bp1zxszhcgatnx****.redis.rds.aliyuncs.com:6379"]
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("cluster config on addr nil", func(t *testing.T) {
		var config *ClusterOptions
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "no cluster addr set:jupiter.redis.test.cluster")
				assert.Nil(t, config)
			}
		}()
		config = ClusterConfig("test")
		assert.Equal(t, len(config.Addr), 1)
	})

}

func TestClusterOption_Normal(t *testing.T) {
	var configStr = `
[jupiter.redis.test]
	[jupiter.redis.test.cluster]
		dialTimeout="2s"
		readTimeout="5s"
		idleTimeout="60s"
		enableAccessLog = false
		addr = ["r-bp1zxszhcgatnx****.redis.rds.aliyuncs.com:6379"]
		username = "root"
		password = "xxxxxx"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	t.Run("cluster option", func(t *testing.T) {
		config := ClusterConfig("test")
		assert.Equal(t, config.DialTimeout, time.Second*2)
		assert.Equal(t, config.ReadTimeout, time.Second*5)
		assert.Equal(t, config.IdleTimeout, time.Minute)
		assert.Equal(t, config.MinIdleConns, 20)

		assert.Equal(t, config.Addr[0], "r-bp1zxszhcgatnx****.redis.rds.aliyuncs.com:6379")
		assert.Equal(t, config.Username, "root")
		assert.Equal(t, config.Password, "xxxxxx")
	})

}
