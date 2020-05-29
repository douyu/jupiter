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
	"context"
	"time"

	"github.com/apache/rocketmq-client-go"
	"github.com/apache/rocketmq-client-go/consumer"
	"github.com/apache/rocketmq-client-go/primitive"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
)

// ConsumerConfig consumer config
type ConsumerConfig struct {
	Enable          bool          `json:"enable" toml:"enable"`
	Addr            []string      `json:"addr" toml:"addr"`
	Topic           string        `json:"topic" toml:"topic"`
	Group           string        `json:"group" toml:"group"`
	DialTimeout     time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	RwTimeout       time.Duration `json:"rwTimeout" toml:"rwTimeout"`
	SubExpression   string        `json:"subExpression" toml:"subExpression"`
	Rate            float64       `json:"rate" toml:"rate"`
	Capacity        int64         `json:"capacity" toml:"capacity"`
	WaitMaxDuration time.Duration `json:"waitMaxDuration" toml:"waitMaxDuration"`

	subscribers  map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
	interceptors []primitive.Interceptor
}

// StdPushConsumerConfig ...
func StdPushConsumerConfig(name string) ConsumerConfig {
	return RawConsumerConfig("jupiter.rocketmq." + name + ".consumer")
}

// RawConsumerConfig 返回配置
func RawConsumerConfig(key string) ConsumerConfig {
	var config = DefaultConsumerConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal config", xlog.String("key", key), xlog.Any("config", config))
	}

	return config
}

// DefaultConsumerConfig ...
func DefaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		DialTimeout:  time.Second * 3,
		RwTimeout:    time.Second * 10,
		subscribers:  make(map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)),
		interceptors: make([]primitive.Interceptor, 0),
	}
}

// WithSubscribe ...
func (config *ConsumerConfig) WithSubscribe(topic string, f func(context.Context, *primitive.MessageExt) error) *ConsumerConfig {
	if config.subscribers == nil {
		config.subscribers = make(map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error))
	}

	if _, ok := config.subscribers[topic]; ok {
		xlog.Panic("duplicated subscribe", xlog.String("topic", topic))
	}

	fn := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			err := f(ctx, msg)
			if err != nil {
				xlog.Error("consumer message", xlog.Any("err", err), xlog.Any("msg", msg))
				return consumer.ConsumeRetryLater, err
			}
		}

		return consumer.ConsumeSuccess, nil
	}
	config.subscribers[topic] = fn
	return config
}

// WithInterceptor ...
func (config *ConsumerConfig) WithInterceptor(fs ...primitive.Interceptor) *ConsumerConfig {
	if config.interceptors == nil {
		config.interceptors = make([]primitive.Interceptor, 0)
	}
	config.interceptors = append(config.interceptors, fs...)
	return config
}

// Build ...
func (config *ConsumerConfig) Build() (rocketmq.PushConsumer, error) {
	// 初始化 PushConsumer
	client, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(config.Group),
		consumer.WithNameServer(config.Addr),
		consumer.WithInterceptor(config.interceptors...),
	)

	if err != nil {
		return nil, err
	}

	for topic, fn := range config.subscribers {
		if err := client.Subscribe(topic, consumer.MessageSelector{}, fn); err != nil {
			return client, err
		}
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return client, err
}
