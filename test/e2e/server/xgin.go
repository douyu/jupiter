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
	"github.com/douyu/jupiter/pkg/server/xgin"
	"github.com/douyu/jupiter/test/e2e/framework"
	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("[xgin] e2e test", func() {
	var server *xgin.Server

	ginkgo.BeforeEach(func() {
		server = xgin.DefaultConfig().MustBuild()
		server.GET("/test", func(r *gin.Context) {
			_, _ = r.Writer.WriteString("hello")
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

	ginkgo.DescribeTable("xgin", func(htc framework.HTTPTestCase) {
		framework.RunHTTPTestCase(htc)
	}, ginkgo.Entry("normal case", framework.HTTPTestCase{
		Conf: &resty.Config{
			Addr: "http://localhost:9091",
		},
		Method:       "GET",
		Path:         "/test",
		ExpectStatus: http.StatusOK,
		ExpectBody:   "hello",
	}))
})
