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

package etcdv3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	defaultConfig := DefaultConfig()
	assert.Equal(t, time.Second*5, defaultConfig.ConnectTimeout)
	assert.Equal(t, false, defaultConfig.BasicAuth)
	assert.Equal(t, []string(nil), defaultConfig.Endpoints)
	assert.Equal(t, false, defaultConfig.Secure)
}

func TestConfigSet(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"localhost"}
	assert.Equal(t, []string{"localhost"}, config.Endpoints)
}
