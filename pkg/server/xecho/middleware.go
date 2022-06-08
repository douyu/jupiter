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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
	"runtime"
	"time"

	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xtrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func extractAID(c echo.Context) string {
	return c.Request().Header.Get("AID")
}

// RecoverMiddleware ...
func recoverMiddleware(logger *xlog.Logger, slowQueryThresholdInMilli int64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			var beg = time.Now()
			var fields = make([]xlog.Field, 0, 8)

			defer func() {
				fields = append(fields, zap.Float64("cost", time.Since(beg).Seconds()))
				if rec := recover(); rec != nil {
					switch rec := rec.(type) {
					case error:
						err = rec
					default:
						err = fmt.Errorf("%v", rec)
					}

					stack := make([]byte, 4096)
					length := runtime.Stack(stack, true)
					fields = append(fields, zap.ByteString("stack", stack[:length]))
				}
				fields = append(fields,
					zap.String("method", ctx.Request().Method),
					zap.Int("code", ctx.Response().Status),
					zap.String("host", ctx.Request().Host),
					zap.String("path", ctx.Request().URL.Path),
				)
				if slowQueryThresholdInMilli > 0 {
					if cost := int64(time.Since(beg)) / 1e6; cost > slowQueryThresholdInMilli {
						fields = append(fields, zap.Int64("slow", cost))
					}
				}
				if err != nil {
					fields = append(fields, zap.String("err", err.Error()))
					logger.Error("access", fields...)
					return
				}
				logger.Info("access", fields...)
			}()

			return next(ctx)
		}
	}
}

func metricServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			beg := time.Now()
			err = next(c)
			method := c.Request().Method + "_" + c.Path()
			peer := c.RealIP()
			if aid := extractAID(c); aid != "" {
				peer += "?aid=" + aid
			}
			metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, method, peer)
			metric.ServerHandleCounter.Inc(metric.TypeHTTP, method, peer, http.StatusText(c.Response().Status))
			return err
		}
	}
}

func traceServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		tracer := xtrace.NewTracer(trace.SpanKindServer)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String("http"),
		}
		return func(c echo.Context) (err error) {
			ctx, span := tracer.Start(c.Request().Context(), c.Request().URL.Path, propagation.HeaderCarrier(c.Request().Header), trace.WithAttributes(attrs...))
			span.SetAttributes(
				semconv.HTTPURLKey.String(c.Request().URL.String()),
				semconv.HTTPTargetKey.String(c.Request().URL.Path),
				semconv.HTTPMethodKey.String(c.Request().Method),
				semconv.HTTPUserAgentKey.String(c.Request().UserAgent()),
				semconv.HTTPClientIPKey.String(c.RealIP()),
			)
			c.SetRequest(c.Request().WithContext(ctx))
			defer span.End()
			return next(c)
		}
	}
}
