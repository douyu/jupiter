// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		assert.Equal(t, time.Second*10, config.DialTimeout)
		assert.Equal(t, "127.0.0.1:9091", config.Address)
		assert.Equal(t, false, config.Direct)
		assert.Equal(t, "panic", config.OnDialError)
	})
}
