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

package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	// TODO(gorexlv): add redis ci
	redisConfig := DefaultRedisConfig()
	redisConfig.Addrs = []string{"localhost:6379"}
	redisConfig.Mode = StubMode
	redisClient := redisConfig.Build()
	err := redisClient.Client.Ping().Err()
	if err != nil {
		t.Errorf("redis ping failed:%v", err)
	}
	st := redisClient.Stub().PoolStats()
	t.Logf("running status %+v", st)
	err = redisClient.Close()
	if err != nil {
		t.Errorf("redis close failed:%v", err)
	}
	st = redisClient.Stub().PoolStats()
	t.Logf("close status %+v", st)
}
