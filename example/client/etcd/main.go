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

package main

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
)

func main() {
	eng := &jupiter.Application{}
	err := eng.Startup(
		func() error {
			client := etcdv3.StdConfig("myetcd").Build()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			// 添加数据
			_, err := client.Put(ctx, "/hello", "jupiter")
			if err != nil {
				xlog.Panic(err.Error())
			}

			// 获取数据
			response, err := client.Get(ctx, "/hello", clientv3.WithPrefix())
			if err != nil {
				xlog.Panic(err.Error())
			}

			xlog.Info("get etcd info", xlog.String("key", string(response.Kvs[0].Key)), xlog.String("value", string(response.Kvs[0].Value)))
			return nil
		},
	)
	if err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
}
