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

	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/core/sentinel"
	"github.com/douyu/jupiter/pkg/core/xtrace"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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

func NewTraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := xtrace.NewTracer(trace.SpanKindServer)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		var remote string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}
		operation, mAttrs := xtrace.ParseFullMethod(info.FullMethod)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
		}
		attrs = append(attrs, mAttrs...)
		if p, ok := peer.FromContext(ctx); ok {
			remote = p.Addr.String()
		}
		attrs = append(attrs, xtrace.PeerAttr(remote)...)
		ctx, span := tracer.Start(ctx, operation, xtrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer func() {
			if err != nil {
				span.RecordError(err)
				s, ok := status.FromError(err)
				if ok {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(s.Code())))
				} else {
					span.SetStatus(codes.Error, err.Error())
				}
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()

		ctx = xlog.NewContext(ctx, xlog.Default(), span.SpanContext().TraceID().String())
		ctx = xlog.NewContext(ctx, xlog.Jupiter(), span.SpanContext().TraceID().String())

		return handler(ctx, req)
	}
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func NewTraceStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		tracer := xtrace.NewTracer(trace.SpanKindServer)
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
		}

		var remote string
		md, ok := metadata.FromIncomingContext(ss.Context())
		if ok {
			md = md.Copy()
		} else {
			md = metadata.MD{}
		}

		operation, mAttrs := xtrace.ParseFullMethod(info.FullMethod)

		attrs = append(attrs, mAttrs...)
		if p, ok := peer.FromContext(ss.Context()); ok {
			remote = p.Addr.String()
		}
		attrs = append(attrs, xtrace.PeerAttr(remote)...)

		ctx, span := tracer.Start(ss.Context(), operation, xtrace.MetadataReaderWriter(md), trace.WithAttributes(attrs...))
		defer span.End()

		ctx = xlog.NewContext(ctx, xlog.Default(), span.SpanContext().TraceID().String())
		ctx = xlog.NewContext(ctx, xlog.Jupiter(), span.SpanContext().TraceID().String())

		return handler(srv, contextedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

func defaultStreamServerInterceptor(logger *xlog.Logger, c *Config) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		beg := time.Now()
		fields := make([]xlog.Field, 0, 8)
		event := "normal"
		defer func() {
			if c.SlowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > c.SlowQueryThresholdInMilli {
					event = "slow"
				}
			}

			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, xlog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				xlog.Any("grpc interceptor type", "unary"),
				xlog.FieldMethod(info.FullMethod),
				xlog.FieldCost(time.Since(beg)),
				xlog.FieldEvent(event),
			)

			for key, val := range getPeer(stream.Context()) {
				fields = append(fields, xlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.Error("access", fields...)
				return
			}

			if c.EnableAccessLog {
				logger.Info("access", fields...)
			}
		}()
		return handler(srv, stream)
	}
}

func defaultUnaryServerInterceptor(logger *xlog.Logger, c *Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		beg := time.Now()
		fields := make([]xlog.Field, 0, 8)
		event := "normal"
		defer func() {
			if c.SlowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > c.SlowQueryThresholdInMilli {
					event = "slow"
				}
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, xlog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				xlog.Any("grpc interceptor type", "unary"),
				xlog.FieldMethod(info.FullMethod),
				xlog.FieldCost(time.Since(beg)),
				xlog.FieldEvent(event),
			)

			for key, val := range getPeer(ctx) {
				fields = append(fields, xlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.Error("access", fields...)
				return
			}

			if c.EnableAccessLog {
				logger.Info("access", fields...)
			}
		}()
		return handler(ctx, req)
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

func getPeer(ctx context.Context) map[string]string {
	peerMeta := make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
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
		peerMeta["clientIP"] = clientIP
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta
}

func NewSentinelUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		entry, blockerr := sentinel.Entry(info.FullMethod,
			sentinel.WithResourceType(base.ResTypeRPC),
			sentinel.WithTrafficType(base.Inbound),
		)
		if blockerr != nil {
			return nil, blockerr
		}

		resp, err := handler(ctx, req)
		entry.Exit(sentinel.WithError(err))

		return resp, err
	}
}
