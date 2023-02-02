package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var addr = "localhost:6379"
var addr2 = "localhost:6379"

func Test_Stub(t *testing.T) {
	config := DefaultConfig()
	t.Run("should panic when addr nil", func(t *testing.T) {
		var client *Client
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "no master or slaves for default", r.(string))
				assert.Nil(t, client)
			}
		}()
		client = config.MustSingleton()
		assert.Nil(t, client) // 不会执行到这里
	})
	t.Run("should not panic when dial err", func(t *testing.T) {
		config.Master.Addr = "1.1.1.1"
		config.OnDialError = "error"
		client, err := config.Build()
		assert.NotNil(t, err)
		assert.Nil(t, client.master)

	})
	t.Run("normal start", func(t *testing.T) {
		config.Master.Addr = addr
		config.name = "test"
		client, err := config.Build()
		assert.Nil(t, err)
		assert.NotNil(t, client)
		err = client.CmdOnMaster().Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("TestStdNewRedisStub ping err %v", err)
		}

	})

}
