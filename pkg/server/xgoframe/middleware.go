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
	"net/http"
	"time"

	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gogf/gf/net/ghttp"
	"go.uber.org/zap"
)

// recoverMiddleware ...
func recoverMiddleware(logger *xlog.Logger, slowQueryThresholdInMilli int64) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		defer func() {
			var beg = time.Now()
			var fields = make([]xlog.Field, 0, 8)

			fields = append(fields, zap.Float64("cost", time.Since(beg).Seconds()))

			if slowQueryThresholdInMilli > 0 {
				if cost := int64(time.Since(beg)) / 1e6; cost > slowQueryThresholdInMilli {
					fields = append(fields, zap.Int64("slow", cost))
				}
			}

			fields = append(fields,
				zap.String("method", r.Method),
				zap.Int("code", r.Response.Status),
				zap.Int("size", r.Response.BufferLength()),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("ip", r.GetClientIp()),
				zap.String("remote_addr", r.RemoteAddr),
			)

			logger.Info("access", fields...)
			return
		}()
		r.Middleware.Next()
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
		span, ctx := trace.StartSpanFromContext(
			r.Context(),
			r.Method+" "+r.URL.Path,
			trace.TagComponent("http"),
			trace.TagSpanKind("server"),
			trace.HeaderExtractor(r.Header),
			trace.CustomTag("http.url", r.URL.Path),
			trace.CustomTag("http.method", r.Method),
			trace.CustomTag("peer.ipv4", r.GetClientIp()),
		)
		r.Request = r.WithContext(ctx)
		defer span.Finish()
		r.Middleware.Next()
	}
}
