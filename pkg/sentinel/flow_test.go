package sentinel

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
)

var _ = ginkgo.Describe("sentinel unit test with config", func() {

	ginkgo.BeforeEach(func() {
		conf.Reset()
		conf.LoadFromDataSource(file.NewDataSource("testdata/flow.toml", false), toml.Unmarshal)
		sentinelReqeust.Reset()
		sentinelSuccess.Reset()
		sentinelExceptionsThrown.Reset()
		sentinelBlocked.Reset()
		sentinelRt.Reset()
	})

	ginkgo.PContext("entry before build", func() {
		ginkgo.It("normal case", func() {

			a, b := Entry("flowrule", api.WithTrafficType(base.Inbound))
			if b != nil {

			} else {
				a.Exit()
			}

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			haveMetric := false
			for _, s := range ss {
				if s.GetName() == "sentinel_request" {
					m := s.GetMetric()[0]
					Expect(m.GetCounter().GetValue()).Should(BeEquivalentTo(1))
					haveMetric = true
					fmt.Println(s)
				}
			}

			Expect(haveMetric).Should(Equal(false))
		})
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

			a, b = Entry("flowrule", api.WithTrafficType(base.Outbound))
			Expect(a).Should(BeNil())
			Expect(b).Should(Not(BeNil()))
			Expect(b.Error()).Should(Equal("SentinelBlockError: BlockTypeFlowControl, message: flow reject check blocked"))
		})
	})
})
