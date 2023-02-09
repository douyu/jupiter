package rocketmq

import (
	"context"
	"sync/atomic"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PullConsumer struct {
	rocketmq.PullConsumer
	name string
	PullConsumerConfig

	subscribers  map[string]func()
	interceptors []primitive.Interceptor
	started      *atomic.Bool
	done         chan struct{}
}

func (conf *PullConsumerConfig) Build() *PullConsumer {
	name := conf.Name

	xlog.Jupiter().Debug("rocketmq's config: ", xlog.String("name", name), xlog.Any("conf", conf))

	cc := &PullConsumer{
		name:               name,
		PullConsumerConfig: *conf,
		subscribers:        make(map[string]func()),
		interceptors:       []primitive.Interceptor{},
		done:               make(chan struct{}, 1),
		started:            new(atomic.Bool),
	}
	cc.interceptors = append(cc.interceptors,
		consumerMetricInterceptor(),
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

func (cc *PullConsumer) Start() error {
	if cc.started.Load() {
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
		consumer.WithPullInterval(cc.PullInterval),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: cc.AccessKey,
			SecretKey: cc.SecretKey,
		}),
	}
	// 增加广播模式
	if cc.PullConsumerConfig.MessageModel == "BroadCasting" {
		opts = append(opts, consumer.WithConsumerModel(consumer.BroadCasting))
	}
	// 初始化 PushConsumer
	client, err := rocketmq.NewPullConsumer(
		opts...,
	)
	if err != nil {
		xlog.Jupiter().Panic("new pull consumer",
			xlog.FieldName(cc.name),
			xlog.FieldExtMessage(cc.PullConsumerConfig),
			xlog.FieldErr(err),
		)
	}
	cc.PullConsumer = client

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "",
	}
	if cc.PullConsumerConfig.SubExpression != "*" {
		selector.Expression = cc.PullConsumerConfig.SubExpression
	}

	if err := client.Subscribe(cc.Topic, selector); err != nil {
		return err
	}

	if err != nil {
		xlog.Jupiter().Panic("subscribe a topic for consuming",
			xlog.FieldName(cc.name),
			xlog.FieldExtMessage(cc.PullConsumerConfig),
			xlog.FieldErr(err),
		)
	}

	if cc.Enable {
		if err := client.Start(); err != nil {
			xlog.Jupiter().Panic("start pull consumer",
				xlog.FieldName(cc.name),
				xlog.FieldExtMessage(cc.PullConsumerConfig),
				xlog.FieldErr(err),
			)
			return err
		}
		for _, fn := range cc.subscribers {
			fn()
		}
		// 在应用退出的时候，保证注销
		hooks.Register(hooks.Stage_BeforeStop, func() {
			cc.Close()
		})
	}

	cc.started.Store(true)
	return nil
}

// MustStart panics when error found.
func (cc *PullConsumer) MustStart() {
	lo.Must0(cc.Start())
}

func (cc *PullConsumer) Poll(ctx context.Context, f func(context.Context, []*primitive.MessageExt) error) {
	if _, ok := cc.subscribers[cc.Topic]; ok {
		xlog.Jupiter().Panic("duplicated register poll message", zap.String("topic", cc.Topic))
	}

	fn := func() {
		xgo.Go(func() {
			tracer := xtrace.NewTracer(trace.SpanKindConsumer)
			for {
				select {
				case <-cc.done:
					rlog.Info("Poll close message handle.", map[string]interface{}{
						rlog.LogKeyConsumerGroup: cc.Group,
					})
					return
				default:
					pullResult, err := cc.PullConsumer.Poll(ctx, cc.PollTimeout)
					if consumer.IsNoNewMsgError(err) {
						continue
					}
					if err != nil {
						xlog.Jupiter().Error("poll error", xlog.FieldErr(err))
						continue
					}
					if f(ctx, pullResult.GetMsgList()) != nil {
						continue
					}
					if cc.EnableTrace {
						for _, msg := range pullResult.GetMsgList() {
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

					cc.PullConsumer.ACK(context.TODO(), pullResult, consumer.ConsumeSuccess)

				}
			}
		})
	}
	cc.subscribers[cc.Topic] = fn
}

func (cc *PullConsumer) Close() {
	close(cc.done)
	err := cc.Shutdown()
	if err != nil {
		xlog.Jupiter().Warn("pull consumer close fail", zap.Error(err))
	}
}
