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
	"net/http"
	"time"

	"github.com/douyu/jupiter/pkg/client/resty"
	"github.com/douyu/jupiter/pkg/server/xecho"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/douyu/jupiter/test/e2e/framework"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("[xecho] e2e test", func() {
	var server *xecho.Server

	ginkgo.BeforeEach(func() {
		server = xecho.DefaultConfig().MustBuild()
		server.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "hello")
		})

		s := &helloworldv1.FooServer{}
		server.GET("/sayhello/get", xecho.GRPCProxyWrapper(s.SayHello))
		server.POST("/sayhello/post", xecho.GRPCProxyWrapper(s.SayHello))

		go func() {
			err := server.Serve()
			assert.Nil(ginkgo.GinkgoT(), err)
		}()
		time.Sleep(time.Second)
	})

	ginkgo.AfterEach(func() {
		_ = server.Stop()
	})

	ginkgo.DescribeTable("xecho ", func(htc framework.HTTPTestCase) {
		framework.RunHTTPTestCase(htc)
	},
		ginkgo.Entry("normal case", framework.HTTPTestCase{
			Conf: &resty.Config{
				Addr: "http://localhost:9091",
			},
			Method:       "GET",
			Path:         "/",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "hello",
		}),

		ginkgo.Entry("grpc proxy get", framework.HTTPTestCase{
			Conf: &resty.Config{
				Addr: "http://localhost:9091",
			},
			Method:       "GET",
			Query:        "name=bob",
			Path:         "/sayhello/get",
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"error":0,"msg":"","data":{"name":"bob","ageNumber":"0","sex":0,"metadata":{}}}`,
		}),

		ginkgo.Entry("grpc proxy post", framework.HTTPTestCase{
			Conf: &resty.Config{
				Addr: "http://localhost:9091",
			},
			Method: "POST",
			Body:   `{"name":"bob"}`,
			Path:   "/sayhello/post",
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			ExpectStatus: http.StatusOK,
			ExpectBody:   `{"error":0,"msg":"","data":{"name":"bob","ageNumber":"0","sex":0,"metadata":{}}}`,
		}),
	)
})
