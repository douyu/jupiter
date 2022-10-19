package redisgo

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/assert"
)

var addr = "redis:6379"

func Test_Stub(t *testing.T) {
	config := DefaultConfig()
	t.Run("should panic when dial err", func(t *testing.T) {
		var client *redis.Client
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "redis stub start err")
				assert.Nil(t, client)
			}
		}()
		client = config.BuildStub()
		assert.Nil(t, client) // 不会执行到这里
	})
	t.Run("should not panic when dial err", func(t *testing.T) {

		config.OnDialError = "error"
		client := config.BuildStub()
		assert.NotNil(t, client)

	})
	t.Run("normal start", func(t *testing.T) {
		config.Addr = addr
		config.name = "test"
		client := config.BuildStub()
		assert.NotNil(t, client)

		err := client.Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("TestStdNewRedisStub ping err %v", err)
		}

		st := client.PoolStats()
		t.Logf("runing status %+v", st)
	})
}

func Test_Cluster(t *testing.T) {
	config := DefaultConfig()
	var client *redis.ClusterClient
	t.Run("should panic when dial err", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "redis cluster client start err")
			}
		}()
		client = config.BuildCluster()
		assert.Nil(t, client) // 不会执行到这里

	})
	t.Run("should not panic when dial err", func(t *testing.T) {

		config.OnDialError = "error"
		client := config.BuildCluster()
		assert.NotNil(t, client)
	})
	t.Run("normal start", func(t *testing.T) {
		config.Addr = addr
		config.name = "test"
		client := config.BuildCluster()
		assert.NotNil(t, client)

		err := client.Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("TestStdNewRedisStub ping err %v", err)
		}

		st := client.PoolStats()
		t.Logf("runing status %+v", st)
	})
}
