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

	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc"
)

func newGRPCClient(config *Config) *grpc.ClientConn {
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	logger := config.logger.With(
		xlog.FieldMod("client.grpc"),
		xlog.FieldAddr(config.Address),
	)
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName))

	cc, err := grpc.DialContext(ctx, config.Address, dialOptions...)

	if err != nil {
		if config.OnDialError == "panic" {
			logger.Panic("dial grpc server", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err))
		} else {
			logger.Error("dial grpc server", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err))
		}
	}
	logger.Info("start grpc client")
	return cc
}
