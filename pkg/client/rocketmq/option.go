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

package rocketmq

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config config...
type Config struct {
	Addresses    []string            `json:"addr" toml:"addr"`
	PushConsumer *PushConsumerConfig `json:"consumer" toml:"consumer" mapstructure:",squash"`
	PullConsumer *PullConsumerConfig `json:"pullConsumer" toml:"pullConsumer" mapstructure:",squash"`
	Producer     *ProducerConfig     `json:"producer" toml:"producer"`
}

type ConsumerDefaultConfig struct {
	Name          string        `json:"name" toml:"name"`
	Enable        bool          `json:"enable" toml:"enable"`
	Addr          []string      `json:"addr" toml:"addr"`
	Topic         string        `json:"topic" toml:"topic"`
	Group         string        `json:"group" toml:"group"`
	DialTimeout   time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	SubExpression string        `json:"subExpression" toml:"subExpression"`
	// 最大重复消费次数
	Reconsume    int32  `json:"reconsume" toml:"reconsume"`
	AccessKey    string `json:"accessKey" toml:"accessKey"`
	SecretKey    string `json:"secretKey" toml:"secretKey"`
	MessageModel string `json:"messageModel" toml:"messageModel"` // 消费模式,默认clustering
	// client实例名，默认会基于Addr字段生成md5，支持多集群
	InstanceName string `json:"instanceName" toml:"instanceName"`
	// 批量消费的最大消息数量，取值范围：[1, 1024]，默认值为1
	ConsumeMessageBatchMaxSize int `json:"consumeMessageBatchMaxSize" toml:"consumeMessageBatchMaxSize"`
	// 每批次从broker拉取消息的最大个数，取值范围：[1, 1024]，默认值为32
	PullBatchSize int32 `json:"pullBatchSize" toml:"pullBatchSize"`
	// 设置每次消息拉取的时间间隔，push模式最大为65535*time.Millisecond
	PullInterval time.Duration `json:"pullInterval" toml:"pullInterval"`
	// 是否开启trace
	EnableTrace bool `json:"enableTrace" toml:"enableTrace"`
}

// PushConsumerConfig push consumer config
type PushConsumerConfig struct {
	ConsumerDefaultConfig
	RwTimeout       time.Duration `json:"rwTimeout" toml:"rwTimeout"`
	Rate            float64       `json:"rate" toml:"rate"`
	Capacity        int64         `json:"capacity" toml:"capacity"`
	WaitMaxDuration time.Duration `json:"waitMaxDuration" toml:"waitMaxDuration"`
	// 消费消息的协程数，默认为20
	ConsumeGoroutineNums int `json:"consumeGoroutineNums" toml:"consumeGoroutineNums"`
}

// PullConsumerConfig pull consumer config
type PullConsumerConfig struct {
	ConsumerDefaultConfig
	// 持久化offset间隔
	RefreshPersistOffsetDuration time.Duration `json:"refreshPersistOffsetDuration" toml:"refreshPersistOffsetDuration"`
	PollTimeout                  time.Duration `json:"pollTimeout" toml:"pollTimeout"`
}

// ProducerConfig producer config
type ProducerConfig struct {
	Name        string        `json:"name" toml:"name"`
	Addr        []string      `json:"addr" toml:"addr"`
	Topic       string        `json:"topic" toml:"topic"`
	Group       string        `json:"group" toml:"group"`
	Retry       int           `json:"retry" toml:"retry"`
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	RwTimeout   time.Duration `json:"rwTimeout" toml:"rwTimeout"`
	AccessKey   string        `json:"accessKey" toml:"accessKey"`
	SecretKey   string        `json:"secretKey" toml:"secretKey"`
	// client实例名，默认会基于Addr字段生成md5，支持多集群
	InstanceName string `json:"instanceName" toml:"instanceName"`
	EnableTrace  bool   `json:"enableTrace" toml:"enableTrace"`
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Addresses: make([]string, 0),
		Producer: &ProducerConfig{
			Retry:       5,
			DialTimeout: time.Second * 3,
			RwTimeout:   0,
			EnableTrace: true,
		},
		PushConsumer: &PushConsumerConfig{
			ConsumerDefaultConfig: ConsumerDefaultConfig{
				Enable:       true,
				DialTimeout:  time.Second * 3,
				Reconsume:    16,
				EnableTrace:  true,
				MessageModel: "Clustering",
			},
			RwTimeout:            time.Second * 10,
			WaitMaxDuration:      60 * time.Second,
			ConsumeGoroutineNums: 20,
		},
		PullConsumer: &PullConsumerConfig{
			ConsumerDefaultConfig: ConsumerDefaultConfig{
				Enable:       true,
				DialTimeout:  time.Second * 3,
				Reconsume:    16,
				EnableTrace:  true,
				MessageModel: "Clustering",
			},
			RefreshPersistOffsetDuration: time.Second * 5,
			PollTimeout:                  time.Second * 5,
		},
	}
}

// StdPushConsumerConfig ...
func StdPushConsumerConfig(name string) *PushConsumerConfig {
	return RawPushConsumerConfig(constant.ConfigKey("rocketmq." + name))
}

// StdPullConsumerConfig ...
func StdPullConsumerConfig(name string) *PullConsumerConfig {
	return RawPullConsumerConfig(constant.ConfigKey("rocketmq." + name))
}

// StdProducerConfig ...
func StdProducerConfig(name string) *ProducerConfig {
	return RawProducerConfig(constant.ConfigKey("rocketmq." + name))
}

// RawPushConsumerConfig 返push consume回配置
// nolint:dupl
func RawPushConsumerConfig(name string) *PushConsumerConfig {
	var defaultConfig = DefaultConfig()
	var pushConsumerConfig = defaultConfig.PushConsumer
	if err := conf.UnmarshalKey(name, &defaultConfig, conf.TagName("toml")); err != nil ||
		(len(pushConsumerConfig.Addr) == 0 && len(defaultConfig.Addresses) == 0) ||
		len(pushConsumerConfig.Topic) == 0 {
		xlog.Jupiter().Panic("pushConsumerConfig fail", xlog.FieldErr(err), xlog.String("key", name), xlog.Any("config", pushConsumerConfig))
	}
	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(pushConsumerConfig.Addr) == 0 {
		pushConsumerConfig.Addr = defaultConfig.Addresses
	}

	pushConsumerConfig.Name = name
	pushConsumerConfig.Addr = compatible(pushConsumerConfig.Addr)

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if pushConsumerConfig.InstanceName == "" {
		pushConsumerConfig.InstanceName = fmt.Sprintf("%x@%d", md5.Sum([]byte(strings.Join(pushConsumerConfig.Addr, ","))), os.Getpid())
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(name, pushConsumerConfig)
	}
	return pushConsumerConfig
}

// RawPullConsumerConfig 返回pull consume配置
// nolint:dupl
func RawPullConsumerConfig(name string) *PullConsumerConfig {
	var defaultConfig = DefaultConfig()
	var pullConsumerConfig = defaultConfig.PullConsumer
	if err := conf.UnmarshalKey(name, &defaultConfig, conf.TagName("toml")); err != nil ||
		(len(pullConsumerConfig.Addr) == 0 && len(defaultConfig.Addresses) == 0) ||
		len(pullConsumerConfig.Topic) == 0 {
		xlog.Jupiter().Panic("PullConsumerConfig fail", xlog.FieldErr(err), xlog.String("key", name), xlog.Any("config", pullConsumerConfig))
	}

	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(pullConsumerConfig.Addr) == 0 {
		pullConsumerConfig.Addr = defaultConfig.Addresses
	}

	pullConsumerConfig.Name = name
	pullConsumerConfig.Addr = compatible(pullConsumerConfig.Addr)

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if pullConsumerConfig.InstanceName == "" {
		pullConsumerConfig.InstanceName = fmt.Sprintf("%x@%d", md5.Sum([]byte(strings.Join(pullConsumerConfig.Addr, ","))), os.Getpid())
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(name, pullConsumerConfig)
	}
	return pullConsumerConfig
}

// RawProducerConfig 返回produce配置
// nolint:dupl
func RawProducerConfig(name string) *ProducerConfig {
	var defaultConfig = DefaultConfig()
	var producerConfig = defaultConfig.Producer
	if err := conf.UnmarshalKey(name, &defaultConfig, conf.TagName("toml")); err != nil ||
		(len(producerConfig.Addr) == 0 && len(defaultConfig.Addresses) == 0) ||
		len(producerConfig.Topic) == 0 {
		xlog.Jupiter().Panic("RawProducerConfig fail", xlog.FieldErr(err), xlog.String("key", name), xlog.Any("config", producerConfig))
	}

	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(producerConfig.Addr) == 0 {
		producerConfig.Addr = defaultConfig.Addresses
	}
	producerConfig.Name = name
	producerConfig.Addr = compatible(producerConfig.Addr)

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if producerConfig.InstanceName == "" {
		producerConfig.InstanceName = fmt.Sprintf("%x@%d", md5.Sum([]byte(strings.Join(producerConfig.Addr, ","))), os.Getpid())
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(name, producerConfig)
	}
	return producerConfig
}

func compatible(addr []string) []string {
	for ind, a := range addr {
		if !strings.HasPrefix(a, "http") {
			addr[ind] = "http://" + a
		}
	}
	return addr
}
