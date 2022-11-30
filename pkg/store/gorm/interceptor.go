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

package gorm

import (
	"errors"
	"time"

	"github.com/alibaba/sentinel-golang/core/base"
	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/core/sentinel"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type Handler func(*gorm.DB)
type Interceptor func(dsn *DSN, op string, options *Config, next Handler) Handler

var errSlowCommand = errors.New("mysql slow command")

func metricInterceptor() Interceptor {
	return func(dsn *DSN, op string, options *Config, next Handler) Handler {
		return func(scope *gorm.DB) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.Error != nil {
				prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr, getStatement(scope.Error.Error())).Inc()

				if scope.Error != gorm.ErrRecordNotFound {
					xlog.Jupiter().Error("mysql err", xlog.FieldErr(scope.Error), xlog.FieldName(dsn.DBName+"."+scope.Name()), xlog.FieldMethod(op))
				} else {
					xlog.Jupiter().Warn("record not found", xlog.FieldErr(scope.Error), xlog.FieldName(dsn.DBName+"."+scope.Name()), xlog.FieldMethod(op))
				}
			} else {
				prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr, "OK").Inc()
			}

			prome.LibHandleHistogram.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr).Observe(cost.Seconds())

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				xlog.Jupiter().Error(
					"slow",
					xlog.FieldErr(errSlowCommand),
					xlog.FieldMethod(op),
					xlog.FieldExtMessage(logSQL(scope.Statement.SQL.String(), scope.Statement.Vars, options.DetailSQL)),
					xlog.FieldAddr(dsn.Addr),
					xlog.FieldName(dsn.DBName+"."+scope.Name()),
					xlog.FieldCost(cost),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}

func traceInterceptor() Interceptor {
	tracer := xtrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("gorm"),
	}

	return func(dsn *DSN, op string, options *Config, next Handler) Handler {
		return func(scope *gorm.DB) {

			if ctx := scope.Statement.Context; ctx != nil {
				md := metadata.New(nil)

				_, span := tracer.Start(ctx, op, propagation.HeaderCarrier(md), trace.WithAttributes(attrs...))

				span.SetAttributes(semconv.DBNameKey.String(dsn.DBName))
				span.SetAttributes(semconv.DBConnectionStringKey.String(dsn.Addr))
				span.SetAttributes(semconv.DBUserKey.String(dsn.User))
				span.SetAttributes(semconv.DBStatementKey.String(
					logSQL(scope.Statement.SQL.String(), scope.Statement.Vars, options.DetailSQL)))

				defer span.End()

				next(scope)

				return
			}

			next(scope)

		}
	}
}

func sentinelInterceptor() Interceptor {
	return func(dsn *DSN, op string, options *Config, next Handler) Handler {
		return func(scope *gorm.DB) {
			entry, blockerr := sentinel.Entry(dsn.Addr,
				sentinel.WithResourceType(base.ResTypeDBSQL),
				sentinel.WithTrafficType(base.Outbound),
			)
			if blockerr != nil {
				_ = scope.AddError(blockerr)

				return
			}

			next(scope)

			entry.Exit(sentinel.WithError(scope.Error))
		}
	}
}
