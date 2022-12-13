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
	"github.com/douyu/jupiter/pkg/core/tests"
	"github.com/douyu/jupiter/pkg/server/xecho"
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
		go func() {
			err := server.Serve()
			assert.Nil(ginkgo.GinkgoT(), err)
		}()
		time.Sleep(time.Second)
	})

	ginkgo.AfterEach(func() {
		_ = server.Stop()
	})

	ginkgo.DescribeTable("xecho ", func(htc tests.HTTPTestCase) {
		tests.RunHTTPTestCase(htc)
	}, ginkgo.Entry("normal case", tests.HTTPTestCase{
		Conf: &resty.Config{
			Addr: "http://localhost:9091",
		},
		Method:       "GET",
		Path:         "/",
		ExpectStatus: http.StatusOK,
		ExpectBody:   "hello",
	}))
})
