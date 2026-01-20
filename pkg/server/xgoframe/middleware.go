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

package xgoframe

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gogf/gf/net/ghttp"
	"go.uber.org/zap"
)

// recoveryMiddleware handles panic recovery
func recoveryMiddleware(logger *xlog.Logger) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)

				logger.Error("recovery",
					zap.ByteString("stack", stack[:length]),
					zap.String("method", r.Method),
					zap.Int("code", r.Response.Status),
					zap.String("host", r.Host),
					zap.String("path", r.URL.Path),
					zap.String("ip", r.GetClientIp()),
					zap.Any("err", err),
				)
			}
		}()
		r.Middleware.Next()
	}
}

// slowLogMiddleware logs slow requests
func slowLogMiddleware(logger *xlog.Logger, slowThreshold time.Duration) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		beg := time.Now()
		r.Middleware.Next()
		cost := time.Since(beg)

		if slowThreshold > 0 && cost >= slowThreshold {
			logger.Error("slow",
				zap.String("method", r.Method),
				zap.Int("code", r.Response.Status),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("ip", r.GetClientIp()),
				xlog.FieldCost(cost),
			)
		}
	}
}

func metricServerInterceptor() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		r.Response.CORSDefault()
		beg := time.Now()
		r.Middleware.Next()

		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("AID"))
		metric.ServerHandleCounter.Inc(metric.TypeHTTP, r.Method+"."+r.URL.Path, r.Header.Get("AID"), http.StatusText(r.Response.Status))
	}
}
func traceServerInterceptor() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		r.Request = r.WithContext(r.Context())
		r.Middleware.Next()
	}
}
