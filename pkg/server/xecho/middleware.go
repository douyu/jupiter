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

package xecho

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/core/sentinel"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func extractAID(c echo.Context) string {
	return c.Request().Header.Get("AID")
}

// recoveryMiddleware handles panic recovery
func recoveryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					switch rec := rec.(type) {
					case error:
						err = rec
					default:
						err = fmt.Errorf("%v", rec)
					}

					stack := make([]byte, 4096)
					length := runtime.Stack(stack, false)

					xlog.J(ctx.Request().Context()).Error("recovery",
						zap.ByteString("stack", stack[:length]),
						zap.String("method", ctx.Request().Method),
						zap.Int("code", ctx.Response().Status),
						zap.String("host", ctx.Request().Host),
						zap.String("path", ctx.Request().URL.Path),
						zap.Any("err", err),
					)
				}
			}()

			return next(ctx)
		}
	}
}

// slowLogMiddleware logs slow requests
func slowLogMiddleware(slowThreshold time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			beg := time.Now()
			err = next(ctx)
			cost := time.Since(beg)

			if slowThreshold > 0 && cost >= slowThreshold {
				xlog.J(ctx.Request().Context()).Error("slow",
					zap.String("method", ctx.Request().Method),
					zap.Int("code", ctx.Response().Status),
					zap.String("host", ctx.Request().Host),
					zap.String("path", ctx.Request().URL.Path),
					xlog.FieldCost(cost),
				)
			}

			return err
		}
	}
}

func metricServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			beg := time.Now()
			err = next(c)
			method := c.Request().Method + "_" + c.Path()
			metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, method, extractAID(c))
			metric.ServerHandleCounter.Inc(metric.TypeHTTP, method, extractAID(c), http.StatusText(c.Response().Status))
			return err
		}
	}
}

func traceServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			tracer := xtrace.NewTracer(trace.SpanKindServer)
			attrs := []attribute.KeyValue{
				semconv.RPCSystemKey.String("http"),
			}

			ctx, span := tracer.Start(c.Request().Context(), c.Request().URL.Path, propagation.HeaderCarrier(c.Request().Header), trace.WithAttributes(attrs...))
			span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest(pkg.Name(), c.Request().URL.Path, c.Request())...)

			ctx = xlog.NewContext(ctx, xlog.Default(), span.SpanContext().TraceID().String())
			ctx = xlog.NewContext(ctx, xlog.Jupiter(), span.SpanContext().TraceID().String())

			c.SetRequest(c.Request().WithContext(ctx))
			defer func() {
				if err != nil {
					span.SetStatus(codes.Error, err.Error())
					span.RecordError(err)
				}

				span.End()
			}()

			return next(c)
		}
	}
}

func sentinelServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			entry, blockerr := sentinel.Entry(c.Request().URL.Host,
				sentinel.WithResourceType(base.ResTypeWeb),
				sentinel.WithTrafficType(base.Inbound),
			)
			if blockerr != nil {
				return blockerr
			}

			err := next(c)
			entry.Exit(sentinel.WithError(err))

			return err
		}
	}
}
