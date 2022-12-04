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
	EnableTrace   bool  `json:"enableTrace" toml:"enableTrace"`
}

/**
// ConsumerConfig consumer config
type ConsumerConfig struct {
	Name            string        `json:"name" toml:"name"`
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
	Shadow          Shadow        `json:"shadow" toml:"shadow"`
	Reconsume       int32         `json:"reconsume" toml:"reconsume"`
	AccessKey       string        `json:"accessKey" toml:"accessKey"`
	SecretKey       string        `json:"secretKey" toml:"secretKey"`
	MessageModel    string        `json:"messageModel" toml:"messageModel"` // 消费模式,默认clustering
	// client实例名，默认会基于Addr字段生成md5，支持多集群
	InstanceName string `json:"instanceName" toml:"instanceName"`
	EnableTrace  bool   `json:"enableTrace" toml:"enableTrace"`
	// 批量消费的最大消息数量，取值范围：[1, 1024]，默认值为1
	ConsumeMessageBatchMaxSize int `json:"consumeMessageBatchMaxSize" toml:"consumeMessageBatchMaxSize"`
	// 每批次从broker拉取消息的最大个数，取值范围：[1, 1024]，默认值为32
	PullBatchSize int32 `json:"pullBatchSize" toml:"pullBatchSize"`
}
*/

// PushConsumerConfig push consumer config
type PushConsumerConfig struct {
	ConsumerDefaultConfig
	RwTimeout       time.Duration `json:"rwTimeout" toml:"rwTimeout"`
	Rate            float64       `json:"rate" toml:"rate"`
	Capacity        int64         `json:"capacity" toml:"capacity"`
	WaitMaxDuration time.Duration `json:"waitMaxDuration" toml:"waitMaxDuration"`
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
			Retry:       3,
			DialTimeout: time.Second * 3,
			RwTimeout:   0,
			EnableTrace: true,
		},
		PushConsumer: &PushConsumerConfig{
			ConsumerDefaultConfig: ConsumerDefaultConfig{
				DialTimeout:  time.Second * 3,
				Reconsume:    16,
				EnableTrace:  true,
				MessageModel: "Clustering",
			},
			RwTimeout:       time.Second * 10,
			WaitMaxDuration: 60 * time.Second,
		},
		PullConsumer: &PullConsumerConfig{
			ConsumerDefaultConfig: ConsumerDefaultConfig{
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
	rc := RawConfig(constant.ConfigKey("rocketmq." + name))
	pc := rc.PushConsumer
	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(rc.PushConsumer.Addr) == 0 {
		pc.Addr = rc.Addresses
	}

	pc.Name = name

	for ind, addr := range pc.Addr {
		if strings.HasPrefix(addr, "http") {
			pc.Addr[ind] = addr
		} else {
			pc.Addr[ind] = "http://" + addr
		}
	}

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if pc.InstanceName == "" {
		pc.InstanceName = fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(pc.Addr, ","))))
	}

	return pc
}

// StdPullConsumerConfig ...
func StdPullConsumerConfig(name string) *PullConsumerConfig {

	rc := RawConfig(constant.ConfigKey("rocketmq." + name))
	pc := rc.PullConsumer
	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(pc.Addr) == 0 {
		pc.Addr = rc.Addresses
	}

	pc.Name = name

	for ind, addr := range pc.Addr {
		if strings.HasPrefix(addr, "http") {
			pc.Addr[ind] = addr
		} else {
			pc.Addr[ind] = "http://" + addr
		}
	}

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if pc.InstanceName == "" {
		pc.InstanceName = fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(pc.Addr, ","))))
	}
	return pc
}

// StdProducerConfig ...
func StdProducerConfig(name string) *ProducerConfig {

	rc := RawConfig(constant.ConfigKey("rocketmq." + name))
	pc := rc.Producer
	// 兼容rocket_client_mq变更，addr需要携带shceme
	if len(pc.Addr) == 0 {
		pc.Addr = rc.Addresses
	}

	pc.Name = name
	for ind, addr := range pc.Addr {
		if strings.HasPrefix(addr, "http") {
			pc.Addr[ind] = addr
		} else {
			pc.Addr[ind] = "http://" + addr
		}
	}

	// 这里根据mq集群地址的md5，生成默认InstanceName
	// 实现自动支持多集群，解决官方库默认不支持多集群消费的问题
	if pc.InstanceName == "" {
		pc.InstanceName = fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(pc.Addr, ","))))
	}

	return pc
}

// RawConfig 返回配置
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("unmarshal config", xlog.FieldErr(err), xlog.String("key", key), xlog.Any("config", config))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, config)
	}
	return config
}
