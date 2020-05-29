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

package rocketmq

import (
	"time"

	"github.com/apache/rocketmq-client-go"
	"github.com/apache/rocketmq-client-go/primitive"
	"github.com/apache/rocketmq-client-go/producer"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
)

// ProducerConfig producer config
type ProducerConfig struct {
	Addr        []string      `json:"addr" toml:"addr"`
	Topic       string        `json:"topic" toml:"topic"`
	Group       string        `json:"group" toml:"group"`
	Retry       int           `json:"retry" toml:"retry"`
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	RwTimeout   time.Duration `json:"rwTimeout" toml:"rwTimeout"`

	interceptors []primitive.Interceptor
}

// StdProducerConfig ...
func StdProducerConfig(name string) ProducerConfig {
	return RawProducerConfig("jupiter.rocketmq." + name + ".producer")
}

// RawProducerConfig ...
func RawProducerConfig(key string) ProducerConfig {
	var config = DefaultProducerConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal config", xlog.String("key", key))
	}

	return config
}

// DefaultProducerConfig ...
func DefaultProducerConfig() ProducerConfig {
	return ProducerConfig{
		Retry:        3,
		DialTimeout:  time.Second * 3,
		RwTimeout:    0,
		interceptors: make([]primitive.Interceptor, 0),
	}
}

// Build ...
func (config ProducerConfig) Build() (rocketmq.Producer, error) {
	// 兼容配置
	client, err := rocketmq.NewProducer(
		producer.WithNameServer(config.Addr),
		producer.WithRetry(config.Retry),
		producer.WithInterceptor(),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return client, err
}

// WithInterceptor ...
func (config *ProducerConfig) WithInterceptor(fs ...primitive.Interceptor) *ProducerConfig {
	if config.interceptors == nil {
		config.interceptors = make([]primitive.Interceptor, 0)
	}
	config.interceptors = append(config.interceptors, fs...)
	return config
}
