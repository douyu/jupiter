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

	"github.com/douyu/jupiter/pkg/client/grpc/resolver"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConn = grpc.ClientConn

func init() {
	// conf.OnLoaded(func(c *conf.Configuration) {
	// xgrpclog.SetLogger(xlog.Jupiter().With(xlog.FieldMod("grpc")))
	// })
}

func newGRPCClient(config *Config) *grpc.ClientConn {
	var ctx = context.Background()

	dialOptions := getDialOptions(config)

	logger := config.logger.With(
		xlog.FieldMod(ecode.ModClientGrpc),
		xlog.FieldAddr(config.Addr),
	)
	// 默认使用block连接，失败后fallback到异步连接
	if config.DialTimeout > time.Duration(0) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)

		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, config.Addr, append(dialOptions, grpc.WithBlock())...)
	if err != nil {
		xlog.Jupiter().Error("dial grpc server failed, connect without block",
			xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err))

		conn, err = grpc.DialContext(context.Background(), config.Addr, dialOptions...)
		if err != nil {
			logger.Error("connect without block failed",
				xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err))
		}
	}

	xlog.Jupiter().Info("start grpc client")

	return conn
}

func getDialOptions(config *Config) []grpc.DialOption {
	dialOptions := config.dialOptions

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolver.NewEtcdBuilder("etcd", config.RegistryConfig)),
	)

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.BalancerName)
	dialOptions = append(dialOptions, grpc.WithDefaultServiceConfig(svcCfg))

	return dialOptions
}
