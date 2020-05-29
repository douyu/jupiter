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
	"time"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/metric"
	"google.golang.org/grpc"
)

// metric统计
func metricUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 收敛err错误，将err过滤后，可以知道err是否为系统错误码
		codes := ecode.ExtractCodes(err)

		// 只记录系统级别的详细错误码
		metric.ClientMetricsHandler.GetHandlerCounter().
			WithLabelValues(metric.TypeServerUnary, name, method, cc.Target(), codes.GetMessage()).Inc()
		metric.ClientMetricsHandler.GetHandlerHistogram().
			WithLabelValues(metric.TypeServerUnary, name, method, cc.Target()).Observe(time.Since(beg).Seconds())
		return err
	}
}

func metricStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		beg := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		// 暂时用默认的grpc的默认err收敛
		codes := ecode.ExtractCodes(err)
		metric.ClientMetricsHandler.GetHandlerCounter().
			WithLabelValues(metric.TypeServerStream, name, method, cc.Target(), codes.GetMessage()).Inc()
		metric.ClientMetricsHandler.GetHandlerHistogram().
			WithLabelValues(metric.TypeServerStream, name, method, cc.Target()).Observe(time.Since(beg).Seconds())
		return clientStream, err
	}
}
