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
	"runtime"
	"time"

	"github.com/douyu/jupiter/pkg/metric"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Config) extractAID(c echo.Context) string {
	return c.Request().Header.Get("AID")
}

// RecoverMiddleware ...
func (invoker *Config) RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			var beg = time.Now()
			var trace = &xlog.Tracer{}
			xlog.InjectTraceMD(ctx, trace)

			defer func() {
				trace.Info(zap.Float64("cost", time.Since(beg).Seconds()))
				if rec := recover(); rec != nil {
					switch rec := rec.(type) {
					case error:
						err = rec
					default:
						err = fmt.Errorf("%v", rec)
					}

					stack := make([]byte, 4096)
					length := runtime.Stack(stack, true)
					trace.Error(zap.ByteString("stack", stack[:length]))
				}
				if err != nil {
					trace.Error(zap.String("err", err.Error()))
				}
				trace.Flush("access", invoker.logger)
			}()

			err = next(ctx)
			return err
		}
	}
}

// AccessLogger ...
func (invoker *Config) AccessLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			err = next(ctx)
			if trace, ok := xlog.ExtractTraceMD(ctx); ok {
				trace.Info(zap.String("method", ctx.Request().Method))
				trace.Info(zap.Int("code", ctx.Response().Status))
				trace.Info(zap.String("host", ctx.Request().Host))
				trace.Info(zap.String("path", ctx.Request().URL.Path))

				if cost := int64(time.Since(trace.BeginTime)) / 1e6; cost > 500 {
					trace.Warn(zap.Int64("slow", cost))
				}
			}
			return err
		}
	}
}

// PrometheusServerInterceptor ...
func (invoker *Config) PrometheusServerInterceptor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			beg := time.Now()
			err = next(c)
			metric.ServerMetricsHandler.GetHandlerHistogram().
				WithLabelValues(metric.TypeServerHttp, c.Request().Method+"."+c.Path(), invoker.extractAID(c)).Observe(time.Since(beg).Seconds())
			metric.ServerMetricsHandler.GetHandlerCounter().
				WithLabelValues(metric.TypeServerHttp, c.Request().Method+"."+c.Path(), invoker.extractAID(c), statusText[c.Response().Status]).Inc()
			return err
		}
	}
}
