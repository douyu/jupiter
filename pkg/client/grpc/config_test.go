package grpc

import (
	"bytes"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var configStr = `
[jupiter.client.test]
	balancerName="swr"
	address="127.0.0.1:9091"
	dialTimeout="10s"
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	t.Run("std config", func(t *testing.T) {
		config := StdConfig("test")
		assert.Equal(t, "swr", config.BalancerName)
		assert.Equal(t, time.Second * 10, config.DialTimeout)
		assert.Equal(t, "127.0.0.1:9091", config.Address)
		assert.Equal(t, false, config.Direct)
		assert.Equal(t, "panic", config.OnDialError)
	})
}
