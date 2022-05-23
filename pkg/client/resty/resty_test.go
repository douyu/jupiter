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

package resty

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/xtrace"
	"github.com/douyu/jupiter/pkg/xtrace/jaeger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestE2ESuites(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "resty e2e test cases")
}

var _ = Describe("normal case", func() {
	It("httptest", func() {
		config := DefaultConfig()
		config.Addr = "http://127.0.0.1:9099"
		config.EnableTrace = true

		jaegerConfig := jaeger.DefaultConfig()
		xtrace.SetGlobalTracer(jaegerConfig.Build())
		//config.MustBuild().R().SetContext().Get()
		res, err := config.MustBuild().R().Get("/hello")
		fmt.Println("============》err", err)
		fmt.Println("============》err", res)
		Expect(err).Should(BeNil())
		time.Sleep(10 * time.Second)
		//Expect(res.Status()).Should(Equal("200 OK"))
	})

	/*It("slowlog", func() {
		config := DefaultConfig()
		config.Addr = "http://127.0.0.1:9099"
		// 测试慢日志
		config.SlowThreshold = time.Millisecond
		_, err := config.MustBuild().R().Get("/hello")

		Expect(err).Should(BeNil())
		//Expect(res.Status()).Should(Equal("200 OK"))
	})

	It("on error", func() {
		//config := DefaultConfig()
		//config.Addr = "http://127.0.0.1:9099"
		//jaegerConfig := jaeger.DefaultConfig()
		//// fast flush trace
		//
		//xtrace.SetGlobalTracer(jaegerConfig.Build())
		//
		//_, err := config.MustBuild().R().Get("/hello")
		//Expect(err.(*url.Error).Err).Should(BeEquivalentTo(errors.New("auto redirect is disabled")))
		////Expect(res.Status()).Should(Equal("302 Found"))
		//time.Sleep(100 * time.Millisecond)
	})*/
})
