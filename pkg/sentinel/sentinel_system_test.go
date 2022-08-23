package sentinel

import (
	"github.com/BurntSushi/toml"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/system_metric"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.FDescribe("sentinel unit test with config", func() {

	ginkgo.BeforeEach(func() {
		stdConfig = StdConfig()
		conf.Reset()
		conf.LoadFromDataSource(file.NewDataSource("testdata/system.toml", false), toml.Unmarshal)
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

			a, b := Entry("system", api.WithTrafficType(base.Inbound))
			Expect(b).Should(BeNil())
			a.Exit()

			system_metric.SetSystemCpuUsage(100)

			a, b = Entry("system", api.WithTrafficType(base.Inbound))
			Expect(a).Should(BeNil())
			Expect(b).Should(Not(BeNil()))
			Expect(b.Error()).Should(Equal("SentinelBlockError: BlockTypeSystem, message: system cpu usage check blocked"))
		})
	})
})
