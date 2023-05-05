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
	"errors"
	"net/url"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2ESuites(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "resty e2e test cases")
}

var _ = Describe("normal case", func() {
	It("httptest", func() {
		config := DefaultConfig()
		config.Addr = "http://localhost:8001"
		config.EnableTrace = true
		config.EnableSentinel = true
		res, err := config.MustBuild().R().Get("/get")
		Expect(err).Should(BeNil())
		Expect(res.Status()).Should(Equal("200 OK"))
	})

	It("slowlog", func() {
		config := DefaultConfig()
		config.Addr = "http://localhost:8001"
		config.EnableSentinel = true
		// 测试慢日志
		config.SlowThreshold = time.Millisecond
		res, err := config.MustBuild().R().Get("/get")

		Expect(err).Should(BeNil())
		Expect(res.Status()).Should(Equal("200 OK"))
	})

	It("on error", func() {
		config := DefaultConfig()
		config.Addr = "http://localhost:8001"
		config.EnableSentinel = true

		res, err := config.MustBuild().R().Get("/status/302")
		Expect(err.(*url.Error).Err).Should(BeEquivalentTo(errors.New("auto redirect is disabled")))
		Expect(res.Status()).Should(Equal("302 FOUND"))
		time.Sleep(100 * time.Millisecond)
	})

	It("no addr", func() {
		Expect(func() {
			config := DefaultConfig()
			config.MustBuild().R().Get("")
		}).Should(Panic())
	})

	It("retry", func() {
		config := DefaultConfig()
		config.Addr = "http://localhost:8001"
		config.RetryCount = 1
		config.Timeout = time.Nanosecond
		// 测试慢日志
		config.SlowThreshold = time.Nanosecond
		res, err := config.MustBuild().R().Get("/get")

		Expect(err.Error()).Should(ContainSubstring("context deadline exceeded (Client.Timeout exceeded while awaiting headers)"))
		Expect(res.Status()).Should(Equal(""))
	})
})
