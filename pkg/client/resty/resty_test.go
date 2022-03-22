package resty

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/trace/jaeger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2ESuites(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "resty e2e test cases")
}

var _ = Describe("normal case", func() {
	It("baidu", func() {
		config := DefaultConfig()
		config.Addr = "http://baidu.com"
		config.EnableTrace = true

		jaegerConfig := jaeger.DefaultConfig()
		// fast flush trace
		jaegerConfig.Reporter.BufferFlushInterval = time.Millisecond
		trace.SetGlobalTracer(jaegerConfig.Build())

		res, err := config.MustBuild().R().Get("/")
		Expect(err).Should(BeNil())
		Expect(res.Status()).Should(Equal("200 OK"))
	})

	It("slowlog", func() {
		config := DefaultConfig()
		config.Addr = "http://baidu.com"
		// 测试慢日志
		config.SlowThreshold = time.Millisecond
		res, err := config.MustBuild().R().Get("")

		Expect(err).Should(BeNil())
		Expect(res.Status()).Should(Equal("200 OK"))
	})

	It("on error", func() {
		config := DefaultConfig()
		config.Addr = "http://baidu.com/not found"
		jaegerConfig := jaeger.DefaultConfig()
		// fast flush trace
		jaegerConfig.Reporter.BufferFlushInterval = time.Millisecond
		trace.SetGlobalTracer(jaegerConfig.Build())

		res, err := config.MustBuild().R().Get("")
		Expect(err.(*url.Error).Err).Should(BeEquivalentTo(errors.New("auto redirect is disabled")))
		Expect(res.Status()).Should(Equal("302 Found"))
		time.Sleep(100 * time.Millisecond)
	})

	It("no addr", func() {
		Expect(func() {
			config := DefaultConfig()
			config.MustBuild().R().Get("")
		}).Should(Panic())
	})
})
