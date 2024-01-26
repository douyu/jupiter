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
	"context"
	"errors"
	"runtime/debug"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/juju/ratelimit"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PushConsumer struct {
	rocketmq.PushConsumer
	name string
	PushConsumerConfig

	subscribers  map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
	interceptors []primitive.Interceptor
	bucket       *ratelimit.Bucket
	started      bool
}

func (conf *PushConsumerConfig) Build() *PushConsumer {
	name := conf.Name

	xlog.Jupiter().Debug("rocketmq's config: ", xlog.String("name", name), xlog.Any("conf", conf))

	var bucket *ratelimit.Bucket
	if conf.Rate > 0 && conf.Capacity > 0 {
		bucket = ratelimit.NewBucketWithRate(conf.Rate, conf.Capacity)
	}

	cc := &PushConsumer{
		name:               name,
		PushConsumerConfig: *conf,
		subscribers:        make(map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)),
		interceptors:       []primitive.Interceptor{},
		bucket:             bucket,
	}
	cc.interceptors = append(cc.interceptors,
		consumerMetricInterceptor(),
		consumerSlowInterceptor(cc.Topic, cc.RwTimeout),
		consumerSentinelInterceptor(cc.Addr),
	)

	// 服务启动前先start
	hooks.Register(hooks.Stage_BeforeRun, func() {
		err := cc.Start()
		if err != nil {
			xlog.Jupiter().Panic("rocketmq consumer start fail", zap.Error(err))
		}
	})

	return cc
}

func (cc *PushConsumer) Close() {
	err := cc.Shutdown()
	if err != nil {
		xlog.Jupiter().Warn("consumer close fail", zap.Error(err))
	}
}

func (cc *PushConsumer) WithInterceptor(fs ...primitive.Interceptor) *PushConsumer {
	cc.interceptors = append(cc.interceptors, fs...)
	return cc
}

func (cc *PushConsumer) RegisterSingleMessage(f func(context.Context, *primitive.MessageExt) error) *PushConsumer {
	if _, ok := cc.subscribers[cc.Topic]; ok {
		xlog.Jupiter().Panic("duplicated register single message", zap.String("topic", cc.Topic))
	}

	tracer := xtrace.NewTracer(trace.SpanKindConsumer)
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("rocketmq"),
		semconv.MessagingRocketmqClientGroupKey.String(cc.Group),
		semconv.MessagingRocketmqClientIDKey.String(cc.InstanceName),
		semconv.MessagingRocketmqConsumptionModelKey.String(cc.MessageModel),
	}

	fn := func(ctx context.Context, msgs ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
		// the recover to prevent panic from causing the coroutine to exit when processing msg.
		defer func() {
			if r := recover(); r != nil {
				xlog.Jupiter().Error("consumer message panic", zap.String("stack", string(debug.Stack())))
				result, err = consumer.ConsumeRetryLater, errors.New("consumer message panic")
			}
		}()
		for _, msg := range msgs {
			var (
				span trace.Span
			)

			if cc.EnableTrace {
				carrier := propagation.MapCarrier{}

				for key, value := range msg.GetProperties() {
					carrier[key] = value
				}

				ctx, span = tracer.Start(ctx, msg.Topic, carrier, trace.WithAttributes(attrs...))

				span.SetAttributes(
					semconv.MessagingRocketmqNamespaceKey.String(msg.Topic),
					semconv.MessagingRocketmqMessageTagKey.String(msg.GetTags()),
				)

				ctx = xlog.NewContext(ctx, xlog.Default(), span.SpanContext().TraceID().String())
				ctx = xlog.NewContext(ctx, xlog.Jupiter(), span.SpanContext().TraceID().String())

				defer span.End()
			}

			if cc.bucket != nil {
				if ok := cc.bucket.WaitMaxDuration(1, cc.WaitMaxDuration); !ok {
					xlog.Jupiter().Warn("too many messages, reconsume later", zap.String("body", string(msg.Body)), zap.String("topic", cc.Topic))
					return consumer.ConsumeRetryLater, nil
				}
			}

			if err := f(ctx, msg); err != nil {
				xlog.Jupiter().Error("consumer message", zap.Error(err), zap.String("field", cc.name), zap.Any("ext", msg))
				if cc.EnableTrace && span != nil {
					span.RecordError(err)
				}
				return consumer.ConsumeRetryLater, err
			}
		}

		return consumer.ConsumeSuccess, nil
	}
	cc.subscribers[cc.Topic] = fn
	return cc
}

func (cc *PushConsumer) RegisterBatchMessage(f func(context.Context, ...*primitive.MessageExt) error) *PushConsumer {
	if _, ok := cc.subscribers[cc.Topic]; ok {
		xlog.Jupiter().Panic("duplicated register batch message", zap.String("topic", cc.Topic))
	}

	tracer := xtrace.NewTracer(trace.SpanKindConsumer)

	fn := func(ctx context.Context, msgs ...*primitive.MessageExt) (result consumer.ConsumeResult, err error) {
		// the recover to prevent panic from causing the coroutine to exit when processing msg.
		defer func() {
			if r := recover(); r != nil {
				xlog.Jupiter().Error("consumer message panic", zap.String("stack", string(debug.Stack())))
				result, err = consumer.ConsumeRetryLater, errors.New("consumer message panic")
			}
		}()
		if cc.EnableTrace {
			for _, msg := range msgs {
				var (
					span trace.Span
				)

				carrier := propagation.MapCarrier{}
				attrs := []attribute.KeyValue{
					semconv.MessagingSystemKey.String("rocketmq"),
					semconv.MessagingDestinationKindKey.String(msg.Topic),
				}
				for key, value := range msg.GetProperties() {
					carrier[key] = value
				}
				ctx, span = tracer.Start(ctx, msg.Topic, carrier, trace.WithAttributes(attrs...))
				defer span.End()
			}
		}

		if cc.bucket != nil {
			if ok := cc.bucket.WaitMaxDuration(int64(len(msgs)), cc.WaitMaxDuration); !ok {
				xlog.Jupiter().Warn("too many messages, reconsume later", zap.String("topic", cc.Topic))
				return consumer.ConsumeRetryLater, nil
			}
		}

		if err := f(ctx, msgs...); err != nil {
			xlog.Jupiter().Error("consumer batch message", zap.Error(err), zap.String("field", cc.name))
			return consumer.ConsumeRetryLater, err
		}

		return consumer.ConsumeSuccess, nil
	}
	cc.subscribers[cc.Topic] = fn
	return cc
}

func (cc *PushConsumer) Start() error {
	if cc.started {
		return nil
	}

	var opts = []consumer.Option{
		consumer.WithGroupName(cc.Group),
		consumer.WithInstance(cc.InstanceName),
		consumer.WithNameServer(cc.Addr),
		consumer.WithMaxReconsumeTimes(cc.Reconsume),
		consumer.WithInterceptor(cc.interceptors...),
		consumer.WithConsumeMessageBatchMaxSize(cc.ConsumeMessageBatchMaxSize),
		consumer.WithPullBatchSize(cc.PullBatchSize),
		consumer.WithConsumeGoroutineNums(cc.ConsumeGoroutineNums),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: cc.AccessKey,
			SecretKey: cc.SecretKey,
		}),
	}
	// 增加广播模式
	if cc.PushConsumerConfig.MessageModel == "BroadCasting" {
		opts = append(opts, consumer.WithConsumerModel(consumer.BroadCasting))
	}
	// 初始化 PushConsumer
	client, err := rocketmq.NewPushConsumer(
		opts...,
	)

	if err != nil {
		xlog.Jupiter().Panic("create push consumer panic",
			xlog.FieldName(cc.name),
			xlog.FieldExtMessage(cc.PushConsumerConfig),
			xlog.FieldErr(err),
		)
	}

	cc.PushConsumer = client

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "",
	}
	if cc.PushConsumerConfig.SubExpression != "*" {
		selector.Expression = cc.PushConsumerConfig.SubExpression
	}

	for topic, fn := range cc.subscribers {
		if err := cc.PushConsumer.Subscribe(topic, selector, fn); err != nil {
			return err
		}
	}

	if cc.Enable {
		if err := client.Start(); err != nil {
			xlog.Jupiter().Panic("start push consumer",
				xlog.FieldName(cc.name),
				xlog.FieldExtMessage(cc.PushConsumerConfig),
				xlog.FieldErr(err),
			)
			return err
		}
		// 在应用退出的时候，保证注销
		hooks.Register(hooks.Stage_BeforeStop, func() {
			cc.Close()
		})
	}

	cc.started = true

	return nil
}

// MustStart panics when error found.
func (cc *PushConsumer) MustStart() {
	lo.Must0(cc.Start())
}
