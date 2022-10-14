package redigo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/spf13/cast"

	"github.com/douyu/jupiter/pkg/xtrace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	prome "github.com/douyu/jupiter/pkg/metric"

	"github.com/douyu/jupiter/pkg/util/xdebug"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis/v8"
)

type redigoContextKeyType struct{}

var ctxBegKey = redigoContextKeyType{}

type interceptor struct {
	beforeProcess         func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess          func(ctx context.Context, cmd redis.Cmder) error
	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func newInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	return &interceptor{
		beforeProcess: func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcess: func(ctx context.Context, cmd redis.Cmder) error {
			return nil
		},
		beforeProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			return ctx, nil
		},
		afterProcessPipeline: func(ctx context.Context, cmds []redis.Cmder) error {
			return nil
		},
	}
}

func (i *interceptor) setBeforeProcess(p func(ctx context.Context, cmd redis.Cmder) (context.Context, error)) *interceptor {
	i.beforeProcess = p
	return i
}

func (i *interceptor) setAfterProcess(p func(ctx context.Context, cmd redis.Cmder) error) *interceptor {
	i.afterProcess = p
	return i
}

func (i *interceptor) setBeforeProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)) *interceptor {
	i.beforeProcessPipeline = p
	return i
}

func (i *interceptor) setAfterProcessPipeline(p func(ctx context.Context, cmds []redis.Cmder) error) *interceptor {
	i.afterProcessPipeline = p
	return i
}
func fixedInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	return newInterceptor(compName, config, logger).
		setBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, ctxBegKey, time.Now()), nil
		}).
		setBeforeProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			return context.WithValue(ctx, ctxBegKey, time.Now()), nil
		}).setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

		if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
			logger.Error("slow",
				xlog.FieldErr(errors.New("redis slow command")),
				xlog.FieldName(cmd.Name()),
				xlog.FieldAddr(config.Addr),
				xlog.FieldCost(cost))
		}
		return nil
	}).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

		if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
			logger.Error("slow",
				xlog.FieldErr(errors.New("redis slow command")),
				xlog.FieldName(getCmdsName(cmds)),
				xlog.FieldAddr(config.Addr),
				xlog.FieldCost(cost))
		}
		return nil
	})

}
func debugInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	addr := config.AddrString()

	return newInterceptor(compName, config, logger).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			if !xdebug.IsDevelopmentMode() {
				return cmd.Err()
			}
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			if err != nil {
				log.Println("redigo print err",
					xlog.FieldName(cmd.Name()),
					xlog.FieldAddr(addr),
					xlog.FieldCost(cost),
					xlog.Any("cmd", cmd.Args()),
					xlog.FieldErr(err),
				)

			} else {
				log.Println("redigo print",
					xlog.FieldName(cmd.Name()),
					xlog.FieldAddr(addr),
					xlog.FieldCost(cost),
					xlog.Any("cmd", cmd.Args()),
					xlog.FieldErr(err),
				)
			}
			return err
		}).
		setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
			if !xdebug.IsDevelopmentMode() {
				for _, cmd := range cmds {
					if cmd.Err() != nil {
						return cmd.Err()
					}
				}
				return nil
			}
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

			for _, cmd := range cmds {
				err := cmd.Err()
				if err != nil {
					log.Println("redigo print err",
						xlog.FieldName(cmd.Name()),
						xlog.FieldAddr(addr),
						xlog.FieldCost(cost),
						xlog.Any("cmd", cmd.Args()),
						xlog.FieldErr(err),
					)
					return err
				} else {
					log.Println("redigo print",
						xlog.FieldName(cmd.Name()),
						xlog.FieldAddr(addr),
						xlog.FieldCost(cost),
						xlog.Any("cmd", cmd.Args()),
						xlog.FieldErr(err),
					)
				}
			}
			return nil
		})
}
func metricInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	addr := config.AddrString()

	return newInterceptor(compName, config, logger).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			prome.LibHandleHistogram.WithLabelValues(prome.TypeRedis, compName, cmd.Name(), addr).Observe(cost.Seconds())
			if err != nil {
				if errors.Is(err, redis.Nil) {
					prome.LibHandleCounter.Inc(prome.TypeRedis, compName, cmd.Name(), addr, "Empty")
					return err
				}
				prome.LibHandleCounter.Inc(prome.TypeRedis, compName, cmd.Name(), addr, "Error")
				return err
			}
			prome.LibHandleCounter.Inc(prome.TypeRedis, compName, cmd.Name(), addr, "OK")
			return nil
		}).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
		name := getCmdsName(cmds)
		prome.LibHandleHistogram.WithLabelValues(prome.TypeRedis, compName, name, addr).Observe(cost.Seconds())
		for _, cmd := range cmds {
			if cmd.Err() != nil {
				if errors.Is(cmd.Err(), redis.Nil) {
					prome.LibHandleCounter.Inc(prome.TypeRedis, compName, name, addr, "Empty")
					return cmd.Err()
				}
				prome.LibHandleCounter.Inc(prome.TypeRedis, compName, name, addr, "Error")
				return cmd.Err()
			}
		}
		prome.LibHandleCounter.Inc(prome.TypeRedis, compName, name, addr, "OK")
		return nil
	})
}
func accessInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	return newInterceptor(compName, config, logger).setAfterProcess(
		func(ctx context.Context, cmd redis.Cmder) error {
			var fields = make([]xlog.Field, 0, 15)
			var err = cmd.Err()
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			fields = append(fields, xlog.FieldMod(compName),
				xlog.FieldMethod(cmd.Name()),
				xlog.Any("req", cmd.Args()),
				xlog.Any("res", response(cmd)),
				xlog.FieldCost(cost))

			// 开启了链路，那么就记录链路id
			if config.EnableTraceInterceptor && otel.GetTracerProvider() != nil {
				fields = append(fields, xlog.String("trace_id", xlog.GetTraceID(ctx)))
			}

			// error
			if err != nil {
				fields = append(fields, xlog.FieldErr(err))
				if errors.Is(err, redis.Nil) {
					logger.Warn("access", fields...)
					return err
				}
				logger.Error("access", fields...)
				return err
			}

			fields = append(fields, xlog.FieldEvent("normal"))
			logger.Info("access", fields...)

			return err
		},
	)
}
func traceInterceptor(compName string, config *Config, logger *xlog.Logger) *interceptor {
	tracer := xtrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.NetHostPortKey.String(config.Addr),
		semconv.DBNameKey.Int(config.DB),
		semconv.DBSystemRedis,
	}

	return newInterceptor(compName, config, logger).
		setBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			ctx, span := tracer.Start(ctx, cmd.FullName(), nil, trace.WithAttributes(attrs...))
			span.SetAttributes(
				semconv.DBOperationKey.String(cmd.Name()),
				semconv.DBStatementKey.String(cast.ToString(cmd.Args())),
			)
			return ctx, nil
		}).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			span := trace.SpanFromContext(ctx)
			if err := cmd.Err(); err != nil && err != redis.Nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			span.End()
			return nil
		}).
		setBeforeProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			ctx, span := tracer.Start(ctx, "pipeline", nil, trace.WithAttributes(attrs...))
			span.SetAttributes(
				semconv.DBOperationKey.String(getCmdsName(cmds)),
			)
			return ctx, nil
		}).
		setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
			span := trace.SpanFromContext(ctx)
			for _, cmd := range cmds {
				if err := cmd.Err(); err != nil && err != redis.Nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					span.End()
					return nil
				}
			}

			span.End()
			return nil
		})
}
func response(cmd redis.Cmder) string {
	switch cmd.(type) {
	case *redis.Cmd:
		return fmt.Sprintf("%v", cmd.(*redis.Cmd).Val())
	case *redis.StringCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StringCmd).Val())
	case *redis.StatusCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StatusCmd).Val())
	case *redis.IntCmd:
		return fmt.Sprintf("%v", cmd.(*redis.IntCmd).Val())
	case *redis.DurationCmd:
		return fmt.Sprintf("%v", cmd.(*redis.DurationCmd).Val())
	case *redis.BoolCmd:
		return fmt.Sprintf("%v", cmd.(*redis.BoolCmd).Val())
	case *redis.CommandsInfoCmd:
		return fmt.Sprintf("%v", cmd.(*redis.CommandsInfoCmd).Val())
	case *redis.StringSliceCmd:
		return fmt.Sprintf("%v", cmd.(*redis.StringSliceCmd).Val())
	default:
		return ""
	}
}
func getCmdsName(cmds []redis.Cmder) string {
	cmdNameMap := map[string]bool{}
	cmdName := []string{}
	for _, cmd := range cmds {
		if !cmdNameMap[cmd.Name()] {
			cmdName = append(cmdName, cmd.Name())
			cmdNameMap[cmd.Name()] = true

		}
	}
	return strings.Join(cmdName, "_")
}
