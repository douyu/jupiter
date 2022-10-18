package redisgo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/douyu/jupiter/pkg/util/xstring"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var configStr = `
[jupiter.redisgo]
    [minerva.redisgo.pkmatch]
        [minerva.redisgo.pkmatch.stub]
            maxIdle=10
            maxActive=50
            dialTimeout="2s"
            readTimeout="2s"
            idleTimeout="60s"
            [minerva.redisgo.pkmatch.stub.master]
                addr="redis://127.0.0.1:6379"
            [minerva.redisgo.pkmatch.stub.slaves]
                addr=[
                    "redis://127.0.0.1:6379",
                ]
	`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))

	t.Run("std config", func(t *testing.T) {
		config := StdStubConfig("test")
		fmt.Println(xstring.Json(config))

	})
}
