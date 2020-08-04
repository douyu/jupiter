package xgf

import (
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gogf/gf/net/ghttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// RecoverMiddleware ...
func recoverMiddleware(logger *xlog.Logger, slowQueryThresholdInMilli int64) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		var beg = time.Now()
		var fields = make([]xlog.Field, 0, 8)
		var brokenPipe bool
		defer func() {
			//Latency
			fields = append(fields, zap.Float64("cost", time.Since(beg).Seconds()))
			if slowQueryThresholdInMilli > 0 {
				if cost := int64(time.Since(beg)) / 1e6; cost > slowQueryThresholdInMilli {
					fields = append(fields, zap.Int64("slow", cost))
				}
			}
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				var er = rec.(error)
				//fields = append(fields, zap.ByteString("stack", stack(3)))
				fields = append(fields, zap.String("err", er.Error()))
				logger.Error("access", fields...)
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					return
				}
				r.Response.ClearBuffer()
				return
			}
			// httpRequest, _ := httputil.DumpRequest(c.Request, false)
			// fields = append(fields, zap.ByteString("request", httpRequest))
			fields = append(fields,
				zap.String("method", r.Method),
				zap.Int("code", r.Response.Status),
				zap.Int("size", r.Response.BufferLength()),
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("ip", r.GetClientIp()),
				//zap.String("err", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
			logger.Info("access", fields...)
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
