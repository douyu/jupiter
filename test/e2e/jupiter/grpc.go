// Copyright 2022 Douyu
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

package jupiter

import (
	"context"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter"
	cetcdv3 "github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/client/grpc"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	"github.com/douyu/jupiter/pkg/core/application"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/registry/etcdv3"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xnet"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/douyu/jupiter/test/e2e/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var host, _ = lo.Must2(xnet.GetLocalMainIP())

var _ = ginkgo.Describe("[jupiter] e2e test", ginkgo.Ordered, func() {
	var app *jupiter.Application

	var _ = ginkgo.BeforeAll(func() {
		err := conf.LoadFromDataSource(file.NewDataSource("./config/jupiter.toml", false), toml.Unmarshal)
		assert.NoError(ginkgo.GinkgoT(), err)

		app = jupiter.DefaultApp()
		server := xgrpc.StdConfig("grpc").MustBuild()
		helloworldv1.RegisterGreeterServiceServer(server.Server, new(helloworldv1.FooServer))
		app.Serve(server)
		// app.SetRegistry(etcdv3.DefaultConfig().MustBuild())
		go func(a *application.Application) {
			err := a.Run()
			assert.Nil(ginkgo.GinkgoT(), err)
		}(app)
		time.Sleep(time.Second)
	})

	var _ = ginkgo.AfterAll(func() {
		_ = app.Stop()
		conf.Reset()
	})

	ginkgo.DescribeTable("jupiter grpc sayhello", func(gtc framework.GRPCTestCase) {
		framework.RunGRPCTestCase(gtc)
	},
		ginkgo.Entry("normal case", framework.GRPCTestCase{
			Conf: &grpc.Config{
				Addr: "localhost:9527",
			},
			Method: "/helloworld.v1.GreeterService/SayHello",
			Args: &helloworldv1.SayHelloRequest{
				Name: "jupiter",
			},
			ExpectError:    nil,
			ExpectMetadata: metadata.MD{"content-type": []string{"application/grpc"}},
			ExpectReply:    &helloworldv1.SayHelloResponse{Data: &helloworldv1.SayHelloResponse_Data{Name: "jupiter"}},
		}),
	)

	ginkgo.DescribeTable("jupiter registry", func(tc framework.ETCDTestCase) {
		framework.RunETCDTestCase(tc)
	},
		ginkgo.Entry("normal case", framework.ETCDTestCase{
			Conf: &etcdv3.Config{
				Config: &cetcdv3.Config{
					Endpoints: []string{"http://localhost:2379"},
				},
			},
			DoFn: func(reg registry.Registry) (interface{}, error) {
				res, err := reg.ListServices(context.Background(), "grpc:e2e.test:v1:dev")
				return res, err
			},
			ExpectError: nil,
			ExpectReply: []*server.ServiceInfo{{Address: host + ":9527"}},
		}),
	)
})
