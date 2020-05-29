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
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis"
)

// RedisStub is an unit that handles master and slave.
type RedisStub struct {
	conf   *RedisConfig
	Client *redis.Client
}

func newRedisStub(config *RedisConfig) *RedisStub {
	stub := &RedisStub{
		conf: config,
	}
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	if err := client.Ping().Err(); err != nil {
		switch config.OnDialError {
		case "panic":
			config.logger.Panic("start redis", xlog.Any("err", err), xlog.Any("config", config))
		default:
			config.logger.Error("start redis", xlog.Any("err", err), xlog.Any("config", config))
		}
	}
	stub.Client = client
	return stub
}
