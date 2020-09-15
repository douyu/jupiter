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

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/douyu/jupiter/pkg/defers"
	"github.com/douyu/jupiter/pkg/istats"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

type Producer struct {
	rocketmq.Producer
	name string
	ProducerConfig
	interceptors []primitive.Interceptor
	fInfo        FlowInfo
}

func StdNewProducer(name string) *Producer {
	return StdProducerConfig("configName").Build()
}

func (conf *ProducerConfig) Build() *Producer {
	name := conf.Topic
	if _, ok := _producers.Load(name); ok {
		xlog.Panic("duplicated load", xlog.String("name", name))
	}

	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint("rocketmq's config: "+name, conf)
	}

	cc := &Producer{
		name:           name,
		ProducerConfig: *conf,
		interceptors:   []primitive.Interceptor{},
		fInfo: FlowInfo{
			FlowInfoBase: istats.NewFlowInfoBase(conf.Shadow.Mode),
			Name:         name,
			Addr:         conf.Addr,
			Topic:        conf.Topic,
			Group:        conf.Group,
			GroupType:    "producer",
		},
	}

	cc.interceptors = append(cc.interceptors, producerDefaultInterceptor(cc), producerMDInterceptor(cc), producerShadowInterceptor(cc, conf.Shadow))

	_producers.Store(name, cc)
	return cc
}

func (pc *Producer) Start() error {
	// 兼容配置
	client, err := rocketmq.NewProducer(
		producer.WithNameServer(pc.Addr),
		producer.WithRetry(pc.Retry),
		producer.WithInterceptor(pc.interceptors...),
		producer.WithInstanceName(pc.name),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: pc.AccessKey,
			SecretKey: pc.SecretKey,
		}),
	)
	if err != nil || client == nil {
		xlog.Panic("create producer",
			xlog.FieldName(pc.name),
			xlog.FieldExtMessage(pc.ProducerConfig),
			xlog.Any("error", err),
		)
	}

	if err := client.Start(); err != nil {
		xlog.Panic("start producer",
			xlog.FieldName(pc.name),
			xlog.FieldExtMessage(pc.ProducerConfig),
			xlog.Any("error", err),
		)
	}

	pc.Producer = client
	// 在应用退出的时候，保证注销
	defers.Register(pc.Close)
	return nil
}

func (pc *Producer) WithInterceptor(fs ...primitive.Interceptor) *Producer {
	pc.interceptors = append(pc.interceptors, fs...)
	return pc
}

func (pc *Producer) Close() error {
	err := pc.Shutdown()
	if err != nil {
		xlog.Warn("consumer close fail", xlog.Any("error", err.Error()))
		return err
	}
	return nil
}

// Send rocketmq发送消息
func (pc *Producer) Send(msg []byte) error {
	m := primitive.NewMessage(pc.Topic, msg)
	_, err := pc.SendSync(context.Background(), m)
	if err != nil {
		xlog.Error("send message error", xlog.Any("msg", msg))
		return err
	}
	return nil
}

// SendWithContext 发送消息
func (pc *Producer) SendWithContext(ctx context.Context, msg []byte) error {
	m := primitive.NewMessage(pc.Topic, msg)
	_, err := pc.SendSync(ctx, m)
	if err != nil {
		xlog.Error("send message error", xlog.Any("msg", msg))
		return err
	}
	return nil
}

// SendWithTag rocket mq 发送消息,可以自定义选择 tag
func (pc *Producer) SendWithTag(msg []byte, tag string) error {
	m := primitive.NewMessage(pc.Topic, msg)
	if tag != "" {
		m.WithTag(tag)
	}

	_, err := pc.SendSync(context.Background(), m)
	if err != nil {
		xlog.Error("send message error", xlog.Any("msg", msg))
		return err
	}
	return nil
}

// SendWithResult rocket mq 发送消息,可以自定义选择 tag 及返回结果
func (pc *Producer) SendWithResult(msg []byte, tag string) (*primitive.SendResult, error) {
	m := primitive.NewMessage(pc.Topic, msg)
	if tag != "" {
		m.WithTag(tag)
	}

	res, err := pc.SendSync(context.Background(), m)
	if err != nil {
		xlog.Error("send message error", xlog.Any("msg", msg))
		return res, err
	}
	return res, nil
}

// SendMsg... 自定义消息格式
func (pc *Producer) SendMsg(msg *primitive.Message) (*primitive.SendResult, error) {
	res, err := pc.SendSync(context.Background(), msg)
	if err != nil {
		xlog.Error("send message error", xlog.Any("msg", msg))
		return res, err
	}
	return res, nil
}
