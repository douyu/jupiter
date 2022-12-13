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
	"context"
	"net/http"
	"time"

	"github.com/douyu/jupiter/pkg/client/resty"
	"github.com/douyu/jupiter/pkg/core/tests"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/util/xtest/server/yell"
	"github.com/douyu/jupiter/proto/testproto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("[xgrpcgateway] e2e test", func() {
	var server *xecho.Server

	ginkgo.BeforeEach(func() {
		mux := runtime.NewServeMux()

		_ = testproto.RegisterGreeterServiceHandlerServer(context.Background(),
			mux, new(yell.FooServer))

		server = xecho.DefaultConfig().MustBuild()
		server.Any("/*", echo.WrapHandler(mux))

		go func() {
			err := server.Serve()
			assert.Nil(ginkgo.GinkgoT(), err)
		}()
		time.Sleep(time.Second)
	})

	ginkgo.AfterEach(func() {
		_ = server.Stop()
	})

	ginkgo.DescribeTable("xgrpcgateway", func(htc tests.HTTPTestCase) {
		tests.RunHTTPTestCase(htc)
	},
		ginkgo.Entry("normal case", tests.HTTPTestCase{
			Conf: &resty.Config{
				Addr: "http://localhost:9091",
			},
			Method:       "POST",
			Path:         "/v1/helloworld.Greeter/SayHello",
			Body:         `{"name":"jupiter"}`,
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"error":0,"msg":"","data":{"name":"jupiter"}}`,
		}),
		ginkgo.Entry("404", tests.HTTPTestCase{
			Conf: &resty.Config{
				Addr: "http://localhost:9091",
			},
			Method:       "POST",
			Path:         "/v1/helloworld.Greeter/SayHelloNotFound",
			Body:         `{"name":"jupiter"}`,
			ExpectStatus: http.StatusNotFound,
			ExpectBody:   `{"code":5,"message":"Not Found","details":[]}`,
		}),
	)
})
