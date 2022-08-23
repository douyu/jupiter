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

package sentinel

import (
	"github.com/BurntSushi/toml"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sentinel unit test with config", func() {

	ginkgo.BeforeEach(func() {
		stdConfig = StdConfig()
		conf.Reset()
		conf.LoadFromDataSource(file.NewDataSource("testdata/flow.toml", false), toml.Unmarshal)
		sentinelReqeust.Reset()
		sentinelSuccess.Reset()
		sentinelExceptionsThrown.Reset()
		sentinelBlocked.Reset()
		sentinelRt.Reset()
	})

	ginkgo.Context("enable sentinel and load rules from files", func() {
		ginkgo.BeforeEach(func() {
			stdConfig = StdConfig()
			stdConfig.Enable = true
			stdConfig.Datasource = "files"
			Expect(stdConfig.Build()).Should(BeNil())
		})

		ginkgo.It("normal case", func() {
			a, b := Entry("flowrule", api.WithTrafficType(base.Inbound))
			Expect(b).Should(BeNil())
			a.Exit()

			a, b = Entry("flowrule", api.WithTrafficType(base.Inbound))
			Expect(a).Should(BeNil())
			Expect(b).Should(Not(BeNil()))
			Expect(b.Error()).Should(Equal("SentinelBlockError: BlockTypeFlowControl, message: flow reject check blocked"))
		})
	})
})
