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

package xgrpc

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/douyu/jupiter/pkg/metric"
	"google.golang.org/grpc"
)

func prometheusUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCUnary, info.FullMethod, extractAID(ctx), code.GetMessage())
	return resp, err
}

func prometheusStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	metric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()))
	metric.ServerHandleCounter.Inc(metric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), code.GetMessage())
	return err
}

func traceUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := trace.StartSpanFromContext(
		ctx,
		info.FullMethod,
		trace.FromIncomingContext(ctx),
		trace.TagComponent("gRPC"),
		trace.TagSpanKind("server.unary"),
	)

	defer span.Finish()

	resp, err := handler(ctx, req)

	if err != nil {
		code := codes.Unknown
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		}
		span.SetTag("code", code)
		ext.Error.Set(span, true)
		span.LogFields(trace.String("event", "error"), trace.String("message", err.Error()))
	}
	return resp, err
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	span, ctx := trace.StartSpanFromContext(
		ss.Context(),
		info.FullMethod,
		trace.FromIncomingContext(ss.Context()),
		trace.TagComponent("gRPC"),
		trace.TagSpanKind("server.stream"),
		trace.CustomTag("isServerStream", info.IsServerStream),
	)
	defer span.Finish()

	return handler(srv, contextedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

// RecoveryStreamServerInterceptor recover interceptor for grpc
func (c *Config) RecoveryStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				c.grpcRecoveryWithXlogError("stream", info.FullMethod, rec.(string))
			}
		}()
		return handler(srv, stream)
	}
}

// RecoveryUnaryServerInterceptor  recover interceptor for grpc
func (c *Config) RecoveryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				c.grpcRecoveryWithXlogError("unary", info.FullMethod, rec.(string))
			}
		}()
		return handler(ctx, req)
	}
}

//LoggerStreamServerIntercept loggerInterceptor for grpc
func (c *Config) LoggerStreamServerIntercept() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var trace = xlog.NewTracer()
		defer trace.Flush("stream access logger", c.logger)
		err = handler(srv, ss)
		if err != nil {
			trace.Error(zap.String("err", err.Error()))
		}
		c.grpcLoggerWithTracer(trace, ss.Context(), info.FullMethod)
		return
	}
}

// LoggerUnaryServerIntercept loggerInterceptor for grpc
func (c *Config) LoggerUnaryServerIntercept() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var trace = xlog.NewTracer()
		defer trace.Flush("unary access logger", c.logger)
		resp, err = handler(xlog.NewContext(ctx, *trace), req)
		if err != nil {
			trace.Error(zap.String("err", err.Error()))
		}
		c.grpcLoggerWithTracer(trace, ctx, info.FullMethod)
		return
	}
}

func (c *Config) grpcRecoveryWithXlogError(interceptorType, method, rec string) {
	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, true)]
	c.logger.Error("grpc server recover",
		xlog.Any("grpc interceptor type", interceptorType),
		xlog.FieldStack(stack),
		xlog.FieldMethod(method),
		xlog.FieldErrKind(rec),
	)
}
func (c *Config) grpcLoggerWithTracer(trace *xlog.Tracer, ctx context.Context, method string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			trace.Info(zap.String("aid", strings.Join(val, ";")))
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			ip, err := getClientIP(ctx)
			if err == nil {
				clientIP = ip
			}
		}
		trace.Info(zap.String("ip", clientIP))
		if val, ok := md["client-host"]; ok {
			trace.Info(zap.String("host", strings.Join(val, ";")))
		}
	}
	trace.Info(zap.String("method", method))
	cost := int64(time.Since(trace.BeginTime)) / 1e6
	if cost > 500 {
		trace.Warn(zap.Int64("slow", cost))
	} else {
		trace.Info(zap.Int64("cost", cost))
	}
}
func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}
