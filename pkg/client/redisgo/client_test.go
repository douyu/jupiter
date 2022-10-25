package redisgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var addr = "redis:6379"
var addr2 = "redis:6379"

func Test_Stub(t *testing.T) {
	config := DefaultConfig()
	t.Run("should panic when addr nil", func(t *testing.T) {
		var client *Client
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, r.(string), "redisgo:no master or slaves for default")
				assert.Nil(t, client)
			}
		}()
		client = config.Build()
		assert.Nil(t, client) // 不会执行到这里
	})
	t.Run("should not panic when dial err", func(t *testing.T) {
		config.Master.Addr = "1.1.1.1"
		config.OnDialError = "error"
		client := config.Build()
		assert.NotNil(t, client.master)

	})
	t.Run("normal start", func(t *testing.T) {
		config.Master.Addr = addr
		config.name = "test"
		client := config.Build()
		assert.NotNil(t, client)

		err := client.CmdOnMaster().Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("TestStdNewRedisStub ping err %v", err)
		}

	})
}
