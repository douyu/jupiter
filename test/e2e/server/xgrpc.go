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

package server

import (
	"time"

	"github.com/douyu/jupiter/pkg/client/grpc"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/douyu/jupiter/test/e2e/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var server *xgrpc.Server

var _ = ginkgo.Describe("[grpc] e2e test", func() {
	var _ = ginkgo.BeforeEach(func() {
		server = xgrpc.DefaultConfig().MustBuild()
		helloworldv1.RegisterGreeterServiceServer(server.Server, new(helloworldv1.FooServer))
		go func() {
			err := server.Serve()
			assert.Nil(ginkgo.GinkgoT(), err)
		}()
		time.Sleep(time.Second)
	})

	var _ = ginkgo.AfterEach(func() {
		_ = server.Stop()
	})

	ginkgo.DescribeTable("xgrpc sayhello", func(gtc framework.GRPCTestCase) {
		framework.RunGRPCTestCase(gtc)
	},
		ginkgo.Entry("normal case", framework.GRPCTestCase{
			Conf: &grpc.Config{
				Addr: "localhost:9092",
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

})
