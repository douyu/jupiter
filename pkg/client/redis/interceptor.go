package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/core/sentinel"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
)

type redigoContextKeyType struct{}

var ctxBegKey = redigoContextKeyType{}

type interceptor struct {
	beforeProcess         func(ctx context.Context, cmd redis.Cmder) (context.Context, error)
	afterProcess          func(ctx context.Context, cmd redis.Cmder) error
	beforeProcessPipeline func(ctx context.Context, cmds []redis.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []redis.Cmder) error
}

func (i *interceptor) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return i.beforeProcess(ctx, cmd)
}

func (i *interceptor) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return i.afterProcess(ctx, cmd)
}

func (i *interceptor) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return i.beforeProcessPipeline(ctx, cmds)
}

func (i *interceptor) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return i.afterProcessPipeline(ctx, cmds)
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
func fixedInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {
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
				xlog.FieldAddr(addr),
				xlog.FieldCost(cost))
		}
		return nil
	}).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

		if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
			logger.Error("slow",
				xlog.FieldErr(errors.New("redis slow command")),
				xlog.FieldType("pipeline"),
				xlog.FieldName(getCmdsName(cmds)),
				xlog.FieldAddr(addr),
				xlog.FieldCost(cost))
		}
		return nil
	})

}
func debugInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {

	return newInterceptor(compName, config, logger).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			fmt.Println(xstring.CallerName(6))
			fmt.Printf("[redis ] %s (%s) :\n", addr, cost) // nolint
			if err != nil {
				fmt.Printf(color.RedString("# %s %+v, ERR=(%s)\n\n", cmd.Name(), cmd.Args(), err.Error())) // nolint
			} else {
				fmt.Printf(color.YellowString("# %s %+v: %s\n\n", cmd.Name(), cmd.Args(), response(cmd))) // nolint
			}
			return nil
		}).
		setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			fmt.Println(xstring.CallerName(8))
			fmt.Printf("[redis pipeline] %s (%s) :\n", addr, cost) // nolint
			for _, cmd := range cmds {
				err := cmd.Err()
				if err != nil {
					fmt.Printf(color.RedString("* %s %+v, ERR=<%s>\n", cmd.Name(), cmd.Args(), err.Error())) // nolint
				} else {
					fmt.Printf(color.YellowString("* %s %+v: %s\n", cmd.Name(), cmd.Args(), response(cmd))) // nolint
				}
			}
			fmt.Print("  \n") // nolint
			return nil
		})
}
func metricInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {

	return newInterceptor(compName, config, logger).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			err := cmd.Err()
			name := strings.ToUpper(cmd.Name())
			prome.LibHandleHistogram.WithLabelValues(prome.TypeRedis, name, addr).Observe(cost.Seconds())
			if err != nil {
				if errors.Is(err, redis.Nil) {
					prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "Empty")
				}
				prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "Error")
			}
			prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "OK")
			return nil
		}).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
		names := strings.ToUpper(getCmdsName(cmds))
		prome.LibHandleHistogram.WithLabelValues(prome.TypeRedis, names, addr).Observe(cost.Seconds())
		for _, cmd := range cmds {
			name := strings.ToUpper(cmd.Name())
			if cmd.Err() != nil {
				if errors.Is(cmd.Err(), redis.Nil) {
					prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "Empty")
				}
				prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "Error")
			}
			prome.LibHandleCounter.Inc(prome.TypeRedis, name, addr, "OK")
		}
		return nil
	})
}
func accessInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {
	return newInterceptor(compName, config, logger).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			var fields = make([]xlog.Field, 0, 15)
			var err = cmd.Err()
			cost := time.Since(ctx.Value(ctxBegKey).(time.Time))
			fields = append(fields, xlog.FieldKey(compName),
				xlog.FieldMethod(cmd.Name()),
				xlog.FieldAddr(addr),
				xlog.Any("req", cmd.Args()),
				xlog.FieldCost(cost))

			// error
			if err != nil {
				fields = append(fields, xlog.FieldErr(err))
				if errors.Is(err, redis.Nil) {
					logger.Warn("access", fields...)
					return nil
				}
				logger.Error("access", fields...)
				return nil
			}
			fields = append(fields, xlog.Any("res", response(cmd)))
			logger.Info("access", fields...)

			return nil
		},
		).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		cost := time.Since(ctx.Value(ctxBegKey).(time.Time))

		for _, cmd := range cmds {
			var fields = make([]xlog.Field, 0, 15)
			var err = cmd.Err()
			fields = append(fields, xlog.FieldKey(compName),
				xlog.FieldType("pipeline"),
				xlog.FieldMethod(cmd.Name()),
				xlog.Any("req", cmd.Args()),
				xlog.FieldCost(cost))

			// error
			if err != nil {
				fields = append(fields, xlog.FieldErr(err))
				if errors.Is(err, redis.Nil) {
					logger.Warn("access", fields...)
					continue
				}
				logger.Error("access", fields...)
				continue
			}
			fields = append(fields, xlog.Any("res", response(cmd)))
			logger.Info("access", fields...)

			continue
		}
		return nil
	})
}
func traceInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {
	tracer := xtrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.NetHostPortKey.String(addr),
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
			} else {
				span.SetStatus(codes.Ok, "ok")
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
			span.SetStatus(codes.Ok, "ok")
			span.End()
			return nil
		})
}
func response(cmd redis.Cmder) string {
	switch recv := cmd.(type) {
	case *redis.Cmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.StringCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.StatusCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.IntCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.DurationCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.BoolCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.CommandsInfoCmd:
		return fmt.Sprintf("%v", recv.Val())
	case *redis.StringSliceCmd:
		return fmt.Sprintf("%v", recv.Val())
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

func sentinelInterceptor(compName string, addr string, config *Config, logger *xlog.Logger) *interceptor {
	return newInterceptor(compName, config, logger).
		setBeforeProcess(func(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
			entry, blockerr := sentinel.Entry(addr,
				sentinel.WithResourceType(base.ResTypeCache),
				sentinel.WithTrafficType(base.Outbound),
			)
			if blockerr != nil {
				return ctx, blockerr
			}

			return sentinel.WithContext(ctx, entry), nil
		}).
		setAfterProcess(func(ctx context.Context, cmd redis.Cmder) error {
			if entry := sentinel.FromContext(ctx); entry != nil {
				var err error
				if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
					err = cmd.Err()
				}

				entry.Exit(sentinel.WithError(err))
			}

			return nil
		}).
		setBeforeProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
			entry, blockerr := sentinel.Entry(addr,
				sentinel.WithResourceType(base.ResTypeCache),
				sentinel.WithTrafficType(base.Outbound),
			)
			if blockerr != nil {
				return ctx, blockerr
			}

			return sentinel.WithContext(ctx, entry), nil
		}).setAfterProcessPipeline(func(ctx context.Context, cmds []redis.Cmder) error {
		if entry := sentinel.FromContext(ctx); entry != nil {
			var err error
			for _, cmd := range cmds {
				// skip redis.Nil error
				if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
					err = cmd.Err()

					break
				}
			}

			entry.Exit(sentinel.WithError(err))
		}

		return nil
	})
}
