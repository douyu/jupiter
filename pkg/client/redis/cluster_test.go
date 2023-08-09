package redis

import (
	"bytes"
	"context"
	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cluster(t *testing.T) {
	config := DefaultConfig()
	t.Run("should panic when addr nil", func(t *testing.T) {
		var client *ClusterClient
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "no cluster for default", r.(string))
				assert.Nil(t, client)
			}
		}()
		client = config.MustClusterSingleton()
		assert.NotNil(t, client)
	})
	t.Run("should not panic when dial err", func(t *testing.T) {
		config.Addr = []string{"1.1.1.1"}
		config.OnDialError = "error"
		client, err := config.BuildCluster()
		assert.NotNil(t, err)
		assert.Nil(t, client.cluster)

	})
	t.Run("normal start", func(t *testing.T) {
		config.Addr = []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}
		client, err := config.BuildCluster()
		assert.Nil(t, err)
		assert.NotNil(t, client)
		err = client.Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("Test_Cluster ping err %v", err)
		}
	})

}

func TestClusterConfig(t *testing.T) {
	assert.Equal(t, constant.ConfigKey("redis", "test", "cluster"), "jupiter.redis.test.cluster")

	var configStr = `
[jupiter.redis]
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
		var config *Config
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
