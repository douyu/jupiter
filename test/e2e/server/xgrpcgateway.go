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

	"github.com/douyu/jupiter/pkg/core/tests"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/proto/testproto"
	"github.com/douyu/jupiter/test/e2e/impl"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("[xgrpcgateway] e2e test", func() {
	var server *xecho.Server

	ginkgo.BeforeEach(func() {
		mux := runtime.NewServeMux()

		testproto.RegisterGreeterHandlerServer(context.Background(),
			mux, new(impl.TestProjectImp))

		server = xecho.DefaultConfig().MustBuild()
		server.Any("/*", echo.WrapHandler(mux))

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

	ginkgo.DescribeTable("xgrpcgateway", func(htc tests.HTTPTestCase) {
		tests.RunHTTPTestCase(htc)
	}, ginkgo.Entry("normal case", tests.HTTPTestCase{
		Host:         "http://localhost:9091",
		Method:       "POST",
		Path:         "/v1/helloworld.Greeter/SayHello",
		ExpectStatus: http.StatusOK,
		ExpectBody:   `{"message":"hello", "id64":"0", "id32":0, "idu64":"0", "idu32":0, "name":"", "done":false}`,
	}))
})
