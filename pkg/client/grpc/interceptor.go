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

package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// metric统计
func metricUnaryClientInterceptor(address string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 收敛err错误，将err过滤后，可以知道err是否为系统错误码
		codes := ecode.ExtractCodes(err)

		// 只记录系统级别的详细错误码
		metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, address, method, cc.Target(), codes.GetMessage())
		metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, address, method, cc.Target())
		return err
	}
}

func metricStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		beg := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		// 暂时用默认的grpc的默认err收敛
		codes := ecode.ExtractCodes(err)
		metric.ClientHandleCounter.Inc(metric.TypeGRPCStream, name, method, cc.Target(), codes.GetMessage())
		metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCStream, name, method, cc.Target())
		return clientStream, err
	}
}

func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(prefix), time.Now().Format("04:05.000"), xcolor.Green("Send: "+method+" | "+xstring.Json(req)))
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Red(prefix), time.Now().Format("04:05.000"), xcolor.Red("Erro: "+err.Error()))
		} else {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(prefix), time.Now().Format("04:05.000"), xcolor.Green("Recv: "+xstring.Json(reply)))
		}

		return err
	}
}

func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		span, ctx := trace.StartSpanFromContext(
			ctx,
			method,
			trace.TagSpanKind("client"),
			trace.TagComponent("grpc"),
		)
		defer span.Finish()

		err := invoker(trace.MetadataInjector(ctx, md), method, req, reply, cc, opts...)
		if err != nil {
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("response_code", code)
			ext.Error.Set(span, true)

			span.LogFields(trace.String("event", "error"), trace.String("message", err.Error()))
		}
		return err
	}
}
