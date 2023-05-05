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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
)

var _ = ginkgo.Describe("sentinel unit test with config", func() {

	ginkgo.BeforeEach(func() {
		stdConfig = StdConfig()
		conf.Reset()
		conf.LoadFromDataSource(file.NewDataSource("testdata/circuitbreaker.toml", false), toml.Unmarshal)

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

			a, b := Entry("test")
			if b != nil {

			} else {
				a.Exit()
			}

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			pass := false
			for _, s := range ss {
				if s.GetName() == "sentinel_request" {
					m := s.GetMetric()[0]
					Expect(m.GetCounter().GetValue()).Should(BeEquivalentTo(1))
					pass = true
					fmt.Println(s)
				}
			}

			fmt.Println("debug:", ss)
			Expect(pass).Should(Equal(true))
		})

		ginkgo.It("none exist resource should not block", func() {

			do := func() {
				a, b := Entry("nonexist")
				if b != nil {

				} else {
					a.Exit(base.WithError(errors.New("test error")))
				}
			}

			for i := 0; i < 100; i++ {
				do()
			}

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			// prometheus.WriteToTextfile("test.txt", prometheus.DefaultGatherer)

			pass := false

			for _, s := range ss {
				if s.GetName() == "sentinel_exceptions_thrown" {
					m := s.GetMetric()[0]
					Expect(m.Counter.GetValue()).Should(BeEquivalentTo(100))
					pass = true
				}
			}
			Expect(pass).Should(Equal(true))
		})

		ginkgo.It("timeout case", func() {

			a, b := Entry("timeout")
			if b != nil {
				ginkgo.Fail(b.Error())
			} else {
				time.Sleep(100 * time.Millisecond)
				fmt.Println("pass")
				a.Exit()
				fmt.Println("exit")
			}

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			// prometheus.WriteToTextfile("test.txt", prometheus.DefaultGatherer)
			pass := false

			for _, s := range ss {
				if s.GetName() == "sentinel_rt" {
					m := s.GetMetric()[0]

					Expect(m.Histogram.GetSampleCount()).Should(BeEquivalentTo(1))
					Expect(m.Histogram.GetSampleSum()).Should(BeNumerically(">=", 0.1))
					pass = true
				}
			}

			Expect(pass).Should(Equal(true))

			for i := 0; i < 9; i++ {
				a, _ := Entry("timeout")
				time.Sleep(2 * time.Millisecond)
				a.Exit(base.WithError(nil))
			}

			a, b = Entry("timeout")
			Expect(b).ShouldNot(BeNil())
		})

		ginkgo.It("error case", func() {

			count := 0
			do := func() {
				a, b := Entry("error")
				count++

				if b != nil {
					fmt.Println("[block]", b.Error())
				} else {
					err := errors.New("test error")
					if count >= 2 {
						err = nil
					}

					a.Exit(base.WithError(err))
				}
			}

			do()
			do()
			do()

			// recovery
			time.Sleep(200 * time.Millisecond)
			do()
			do()

			// prometheus.WriteToTextfile("test.txt", prometheus.DefaultGatherer)

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			var pass1, pass2, pass3 bool
			for _, s := range ss {
				if s.GetName() == "sentinel_rt" {
					m := s.GetMetric()[0]
					Expect(m.Histogram.GetSampleCount()).Should(BeEquivalentTo(4))
					pass1 = true
				}

				if s.GetName() == "sentinel_request" {
					m := s.GetMetric()[0]
					fmt.Println(m.GetLabel()[1].GetValue(), "!!!!!", m.GetLabel()[2].GetValue())
					Expect(m.Counter.GetValue()).Should(BeEquivalentTo(5))
					pass2 = true
				}

				if s.GetName() == "sentinel_state" {
					m := s.GetMetric()[0]
					fmt.Println(m)

					Expect(m.Gauge.GetValue()).Should(BeEquivalentTo(circuitbreaker.Closed))
					pass3 = true
				}
			}

			Expect(pass1).Should(Equal(true))
			Expect(pass2).Should(Equal(true))
			Expect(pass3).Should(Equal(true))
		})

	})

	ginkgo.Context("disable sentinel", func() {
		ginkgo.BeforeEach(func() {
			stdConfig = StdConfig()
			stdConfig.Enable = false
			Expect(stdConfig.Build()).Should(BeNil())
		})

		ginkgo.It("no config found in toml", func() {

			a, b := Entry("test")
			if b != nil {

			} else {
				a.Exit()
			}

			// prometheus.WriteToTextfile("test.txt", prometheus.DefaultGatherer)
			Expect(b).Should(BeNil())
			Expect(a).ShouldNot(BeNil())

			ss, err := prometheus.DefaultGatherer.Gather()
			Expect(err).Should(BeNil())

			haveMetric := false
			for _, s := range ss {
				if s.GetName() == "sentinel_request" {
					m := s.GetMetric()[0]
					Expect(m.GetCounter().GetValue()).Should(BeEquivalentTo(1))
					haveMetric = true
				}
			}

			Expect(haveMetric).Should(Equal(false))
		})
	})

	ginkgo.Context("load rules", func() {
		ginkgo.BeforeEach(func() {
			Expect(stdConfig.Build()).Should(BeNil())
			stdConfig.Enable = true
			stdConfig.Datasource = "etcd"
		})

		ginkgo.It("load rules from etcd", func() {

			a, b := Entry("test")
			if b != nil {

			} else {
				a.Exit()
			}

		})
	})

	ginkgo.Context("watch rules", func() {
		var cli *etcdv3.Client

		ginkgo.BeforeEach(func() {
			stdConfig = StdConfig()
			stdConfig.Enable = true
			stdConfig.Datasource = "etcd"

			var err error

			cli, err = etcdv3.RawConfig(stdConfig.EtcdRawKey).Singleton()
			if err != nil {
				ginkgo.Fail("failed to get etcdv3 client")
			}

			_, err = cli.Put(context.Background(),
				"/wsd-sentinel/go/sentinel.test/unknown/local-live/degrade",
				`[{"enable":false,"resource":"test-watch","strategy":2,"retryTimeoutMs":5000,"minRequestAmount":2,"maxAllowedRtMs":100,"statIntervalMs":1000,"statSlidingWindowBucketCount":5,"threshold":0.5}]`)
			Expect(err).Should(BeNil())

			err = stdConfig.Build()
			Expect(err).Should(BeNil())

			time.Sleep(100 * time.Millisecond) //watch操作里每隔一秒拉取一次etcd

			_, err = cli.Put(context.Background(),
				"/wsd-sentinel/go/sentinel.test/unknown/local-live/degrade",
				`[{"enable":true,"resource":"test-watch","strategy":2,"retryTimeoutMs":5000,"minRequestAmount":5,"maxAllowedRtMs":100,"statIntervalMs":1000,"statSlidingWindowBucketCount":5,"threshold":0.5}]`)
			Expect(err).Should(BeNil())

			time.Sleep(100 * time.Millisecond) //watch操作里每隔一秒拉取一次etcd
		})

		ginkgo.AfterEach(func() {
			_, err := cli.Delete(context.Background(),
				"/wsd-sentinel/go/sentinel.test/unknown/local-live/degrade")
			Expect(err).Should(BeNil())
		})

		ginkgo.It("watch rules from etcd", func() {
			time.Sleep(100 * time.Millisecond) //watch操作里每隔一秒拉取一次etcd

			count := 0
			for i := 0; i < 10; i++ {
				a, b := Entry("test-watch")
				if b != nil {
					fmt.Println(count, ":", b.Error())
					count++
				} else {
					fmt.Println(count, ": normal")
					a.Exit(WithError(errors.New("error")))
				}
			}

			Expect(count).Should(BeEquivalentTo(5))
		})
	})
})
