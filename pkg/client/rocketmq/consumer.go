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
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/defers"
	"github.com/douyu/jupiter/pkg/istats"
	"github.com/douyu/jupiter/pkg/xlog"
)

type PushConsumer struct {
	rocketmq.PushConsumer
	name string
	ConsumerConfig

	subscribers  map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
	interceptors []primitive.Interceptor
	fInfo        FlowInfo
}

func (conf *ConsumerConfig) Build() *PushConsumer {
	name := conf.Name
	if _, ok := _consumers.Load(name); ok {
		xlog.Panic("duplicated load", xlog.String("name", name))
	}

	xlog.Debug("rocketmq's config: ", xlog.String("name", name), xlog.Any("conf", conf))

	pc := &PushConsumer{
		name:           name,
		ConsumerConfig: *conf,
		subscribers:    make(map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)),
		interceptors:   []primitive.Interceptor{},
		fInfo: FlowInfo{
			FlowInfoBase: istats.NewFlowInfoBase(conf.Shadow.Mode),
			Name:         name,
			Addr:         conf.Addr,
			Topic:        conf.Topic,
			Group:        conf.Group,
			GroupType:    "consumer",
		},
	}
	pc.interceptors = append(pc.interceptors, pushConsumerDefaultInterceptor(pc), pushConsumerMDInterceptor(pc), pushConsumerShadowInterceptor(pc, conf.Shadow))

	_consumers.Store(name, pc)
	return pc
}

func (cc *PushConsumer) Close() error {
	err := cc.Shutdown()
	if err != nil {
		xlog.Warn("consumer close fail", xlog.Any("error", err.Error()))
		return err
	}
	return nil
}

func (cc *PushConsumer) WithInterceptor(fs ...primitive.Interceptor) *PushConsumer {
	cc.interceptors = append(cc.interceptors, fs...)
	return cc
}

func (cc *PushConsumer) Subscribe(topic string, f func(context.Context, *primitive.MessageExt) error) *PushConsumer {
	if _, ok := cc.subscribers[topic]; ok {
		xlog.Panic("duplicated subscribe", xlog.String("topic", topic))
	}
	fn := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			err := f(ctx, msg)
			if err != nil {
				xlog.Error("consumer message", xlog.String("err", err.Error()), xlog.String("field", cc.name), xlog.Any("ext", msg))
				return consumer.ConsumeRetryLater, err
			}
		}

		return consumer.ConsumeSuccess, nil
	}
	cc.subscribers[topic] = fn
	return cc
}

func (cc *PushConsumer) Start() error {
	// 初始化 PushConsumer
	client, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cc.Group),
		consumer.WithNameServer(cc.Addr),
		consumer.WithMaxReconsumeTimes(cc.Reconsume),
		consumer.WithInterceptor(cc.interceptors...),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: cc.AccessKey,
			SecretKey: cc.SecretKey,
		}),
	)
	cc.PushConsumer = client

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "",
	}
	if cc.ConsumerConfig.SubExpression != "*" {
		selector.Expression = cc.ConsumerConfig.SubExpression
	}

	for topic, fn := range cc.subscribers {
		if err := cc.PushConsumer.Subscribe(topic, selector, fn); err != nil {
			return err
		}
	}

	if err != nil || client == nil {
		xlog.Panic("create consumer",
			xlog.FieldName(cc.name),
			xlog.FieldExtMessage(cc.ConsumerConfig),
			xlog.Any("error", err),
		)
	}

	if cc.Enable {
		if err := client.Start(); err != nil {
			xlog.Panic("start consumer",
				xlog.FieldName(cc.name),
				xlog.FieldExtMessage(cc.ConsumerConfig),
				xlog.Any("error", err),
			)
			return err
		}
		// 在应用退出的时候，保证注销
		defers.Register(cc.Close)
	}

	return nil
}
