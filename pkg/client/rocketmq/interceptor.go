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
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/imeta"
	"github.com/douyu/jupiter/pkg/istats"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

type FlowInfo struct {
	Name      string   `json:"name"`
	Addr      []string `json:"addr"`
	Topic     string   `json:"topic"`
	Group     string   `json:"group"`
	GroupType string   `json:"groupType"` // 类型， consumer 消费者， producer 生产者
	istats.FlowInfoBase
}

func consumeResultStr(result consumer.ConsumeResult) string {
	switch result {
	case consumer.ConsumeSuccess:
		return "success"
	case consumer.ConsumeRetryLater:
		return "retryLater"
	case consumer.Commit:
		return "commit"
	case consumer.Rollback:
		return "rollback"
	case consumer.SuspendCurrentQueueAMoment:
		return "suspendCurrentQueueAMoment"
	default:
		return "unknown"
	}
}

func pushConsumerDefaultInterceptor(pushConsumer *PushConsumer) primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		beg := time.Now()
		msgs := req.([]*primitive.MessageExt)

		err := next(ctx, msgs, reply)
		if reply == nil {
			return err
		}

		holder := reply.(*consumer.ConsumeResultHolder)
		xdebug.PrintObject("consume", map[string]interface{}{
			"err":    err,
			"count":  len(msgs),
			"result": consumeResultStr(holder.ConsumeResult),
		})

		// 消息处理结果统计
		for _, msg := range msgs {
			host := msg.StoreHost
			topic := msg.Topic
			result := consumeResultStr(holder.ConsumeResult)
			if err != nil {
				xlog.Error("push consumer",
					xlog.String("topic", topic),
					xlog.String("host", host),
					xlog.String("result", result),
					xlog.Any("err", err))

			} else {
				xlog.Info("push consumer",
					xlog.String("topic", topic),
					xlog.String("host", host),
					xlog.String("result", result),
				)
			}
			metric.ClientHandleCounter.Inc(metric.TypeRocketMQ, topic, "consume", host, result)
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeRocketMQ, topic, "consume", host)
		}
		if pushConsumer.RwTimeout > time.Duration(0) {
			if time.Since(beg) > pushConsumer.RwTimeout {
				xlog.Error("slow",
					xlog.String("topic", pushConsumer.Topic),
					xlog.String("result", consumeResultStr(holder.ConsumeResult)),
					xlog.Any("cost", time.Since(beg).Seconds()),
				)
			}
		}

		return err
	}
}

func pushConsumerMDInterceptor(pushConsumer *PushConsumer) primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		msgs := req.([]*primitive.MessageExt)
		if len(msgs) > 0 {
			var meta = imeta.New(nil)
			for key, vals := range msgs[0].GetProperties() {
				if strings.HasPrefix(strings.ToLower(key), "x-dy") {
					meta.Set(key, strings.Split(vals, ",")...)
				}
			}
			ctx = imeta.WithContext(ctx, meta)
		}
		err := next(ctx, msgs, reply)
		return err
	}
}

func produceResultStr(result primitive.SendStatus) string {
	switch result {
	case primitive.SendOK:
		return "sendOk"
	case primitive.SendFlushDiskTimeout:
		return "sendFlushDiskTimeout"
	case primitive.SendFlushSlaveTimeout:
		return "sendFlushSlaveTimeout"
	case primitive.SendSlaveNotAvailable:
		return "sendSlaveNotAvailable"
	default:
		return "unknown"
	}
}

func producerDefaultInterceptor(producer *Producer) primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		beg := time.Now()
		realReq := req.(*primitive.Message)
		realReply := reply.(*primitive.SendResult)

		var (
			span opentracing.Span
		)

		if producer.EnableTrace {
			span, ctx = trace.StartSpanFromContext(
				ctx,
				realReq.Topic,
				trace.TagComponent("rocketmq"),
				ext.SpanKindProducer,
			)
			defer span.Finish()

			md := metadata.New(nil)
			ctx = trace.HeaderInjector(ctx, md)
			for k, v := range md {
				realReq.WithProperty(k, strings.Join(v, ","))
			}
		}

		err := next(ctx, realReq, realReply)
		if realReply == nil || realReply.MessageQueue == nil {
			return err
		}

		xdebug.PrintObject("produce", map[string]interface{}{
			"err":     err,
			"message": realReq,
			"result":  realReply.String(),
		})

		if producer.EnableTrace {
			if err != nil {
				ext.Error.Set(span, true)
				ext.LogError(span, err, log.Object("req", realReq), log.Object("res", realReply))
			}
		}

		// 消息处理结果统计
		topic := producer.Topic
		if err != nil {
			xlog.Error("produce",
				xlog.String("topic", topic),
				xlog.String("queue", ""),
				xlog.String("result", realReply.String()),
				xlog.Any("err", err),
			)
			metric.ClientHandleCounter.Inc(metric.TypeRocketMQ, topic, "produce", "unknown", err.Error())
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeRocketMQ, topic, "produce", "unknown")
		} else {
			xlog.Info("produce",
				xlog.String("topic", topic),
				xlog.Any("queue", realReply.MessageQueue),
				xlog.String("result", produceResultStr(realReply.Status)),
			)
			metric.ClientHandleCounter.Inc(metric.TypeRocketMQ, topic, "produce", realReply.MessageQueue.BrokerName, produceResultStr(realReply.Status))
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeRocketMQ, topic, "produce", realReply.MessageQueue.BrokerName)
		}

		if producer.RwTimeout > time.Duration(0) {
			if time.Since(beg) > producer.RwTimeout {
				xlog.Error("slow",
					xlog.String("topic", topic),
					xlog.String("result", realReply.String()),
					xlog.Any("cost", time.Since(beg).Seconds()),
				)
			}
		}

		return err
	}
}

// 统一minerva metadata 传递
func producerMDInterceptor(producer *Producer) primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		if md, ok := imeta.FromContext(ctx); ok {
			realReq := req.(*primitive.Message)
			for k, v := range md {
				realReq.WithProperty(k, strings.Join(v, ","))
			}
		}
		err := next(ctx, req, reply)
		return err
	}
}
