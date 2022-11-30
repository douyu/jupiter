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

	"github.com/douyu/jupiter/pkg/core/tests"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/proto/testproto/v1"
	"github.com/douyu/jupiter/test/e2e/impl"
	"github.com/onsi/ginkgo/v2"
	"google.golang.org/grpc/metadata"
)

var _ = ginkgo.Describe("[grpc] e2e test", func() {
	var server *xgrpc.Server

	ginkgo.BeforeEach(func() {
		server = xgrpc.DefaultConfig().MustBuild()
		testproto.RegisterGreeterServiceServer(server.Server, new(impl.TestProjectImp))
		go func() {
			err := server.Serve()
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(time.Second)
	})

	ginkgo.AfterEach(func() {
		_ = server.Stop()
	})

	ginkgo.DescribeTable("xgrpc sayhello", func(gtc tests.GRPCTestCase) {
		tests.RunGRPCTestCase(gtc)
	}, ginkgo.Entry("normal case", tests.GRPCTestCase{
		Addr:   "localhost:9092",
		Method: "/testproto.v1.GreeterService/SayHello",
		Args: &testproto.SayHelloRequest{
			Name: "jupiter",
		},
		ExpectError:    nil,
		ExpectMetadata: metadata.MD{"content-type": []string{"application/grpc"}},
		ExpectReply:    &testproto.SayHelloResponse{Data: &testproto.SayHelloResponse_Data{Name: "jupiter"}},
	}))

	ginkgo.DescribeTable("xgrpc ", func(gtc tests.GRPCTestCase) {
		tests.RunGRPCTestCase(gtc)
	}, ginkgo.Entry("normal case", tests.GRPCTestCase{
		Addr:   "localhost:9092",
		Method: "/testproto.v1.GreeterService/SayHello",
		Args: &testproto.SayHelloRequest{
			Name: "jupiter",
		},
		ExpectError:    nil,
		ExpectMetadata: metadata.MD{"content-type": []string{"application/grpc"}},
		ExpectReply:    &testproto.SayHelloResponse{Data: &testproto.SayHelloResponse_Data{Name: "jupiter"}},
	}))
})
