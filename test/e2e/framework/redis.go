// Copyright 2022 Douyu
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

package framework

import (
	"context"
	"time"

	"github.com/douyu/jupiter/pkg/client/redis"
	"github.com/imdario/mergo"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type RedisTestCase struct {
	Conf     *redis.Config
	Args     []interface{}
	Timeout  time.Duration
	OnMaster bool

	ExpectError  error
	ExpectResult interface{}
}

// RunRedisTestCase runs a test case against the given handler.
func RunRedisTestCase(rtc RedisTestCase) {
	ginkgoT := ginkgo.GinkgoT()

	if rtc.Timeout == 0 {
		rtc.Timeout = time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), rtc.Timeout)
	defer cancel()

	err := mergo.Merge(rtc.Conf, redis.DefaultConfig())
	assert.Nil(ginkgoT, err)

	client := rtc.Conf.MustBuild()

	stub := client.CmdOnSlave()
	if rtc.OnMaster {
		stub = client.CmdOnMaster()
	}

	result, err := stub.Do(ctx, rtc.Args...).Result()

	assert.Equal(ginkgoT, rtc.ExpectError, err, "error: %s", err)

	if rtc.ExpectResult != nil {
		assert.Equal(ginkgoT, rtc.ExpectResult, result,
			"expected: %s\nactually: %s", rtc.ExpectResult, result)
	}
}
