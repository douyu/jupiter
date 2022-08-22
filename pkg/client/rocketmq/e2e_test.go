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

package rocketmq_test

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/client/rocketmq"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/conf/datasource/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2ESuites(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rocketmq e2e test cases")
}

var _ = Describe("push and consume", func() {

	It("normal case", func() {
		// - NAME_SERVER_ADDRESS=127.0.0.1:9876
		// - BROKER_ADDRESS=127.0.0.1:10911
		// - TOPIC=test
		// - GROUP=testGroup
		conf.LoadFromDataSource(file.NewDataSource("../../../test/testdata/rocketmq/conf/rocketmq.toml", false), toml.Unmarshal)
		consumerClient := rocketmq.StdPushConsumerConfig("example").Build()

		count := int32(0)

		consumerClient.Subscribe(consumerClient.ConsumerConfig.Topic, func(ctx context.Context, ext *primitive.MessageExt) error {
			atomic.AddInt32(&count, 1)
			fmt.Println("msg...", string(ext.Message.Body), string(ext.Message.Topic), string(ext.Message.GetTags()), atomic.LoadInt32(&count))

			return nil
		})
		err := consumerClient.Start()
		Expect(err).Should(BeNil())

		// Eventually(func() int {
		// 	return int(atomic.LoadInt32(&count))
		// }, 10*time.Second, 500*time.Millisecond).Should(Equal(100))

		producerClient := rocketmq.StdProducerConfig("example").Build()

		err = producerClient.Start()
		Expect(err).Should(BeNil())

		for i := 0; i < 10; i++ {
			msg := "d" + strconv.Itoa(i)
			err = producerClient.Send([]byte(msg))
			Expect(err).Should(BeNil())
		}

		for i := 0; i < 10; i++ {
			msg := "a" + strconv.Itoa(i)
			err = producerClient.SendWithTag([]byte(msg), "TagB")
			Expect(err).Should(BeNil())
		}

		Eventually(func() int {
			return int(atomic.LoadInt32(&count))
		}, 5*time.Second, 500*time.Millisecond).Should(Equal(10))

		for i := 0; i < 10; i++ {
			msg := primitive.NewMessage("", []byte("msg"+strconv.Itoa(i)))
			msg = msg.WithTag("TagB")

			err = producerClient.SendWithMsg(context.TODO(), msg)
			Expect(err).Should(BeNil())
		}

		Eventually(func() int {
			return int(atomic.LoadInt32(&count))
		}, 5*time.Second, 500*time.Millisecond).Should(Equal(20))

		consumerClient.Close()
		producerClient.Close()
	})

	It("panic recover", func() {
		conf.LoadFromDataSource(file.NewDataSource("../../../test/testdata/rocketmq/conf/rocketmq.toml", false), toml.Unmarshal)
		consumerClient := rocketmq.StdPushConsumerConfig("example").Build()

		count := int32(0)

		consumerClient.Subscribe(consumerClient.ConsumerConfig.Topic, func(ctx context.Context, ext *primitive.MessageExt) error {
			atomic.AddInt32(&count, 1)
			fmt.Println("msg...", string(ext.Message.Body), string(ext.Message.Topic), string(ext.Message.GetTags()), atomic.LoadInt32(&count))
			panic("test panic")
		})
		err := consumerClient.Start()
		Expect(err).Should(BeNil())

		// Eventually(func() int {
		// 	return int(atomic.LoadInt32(&count))
		// }, 1*time.Second, 500*time.Millisecond).Should(Equal(0))

		producerClient := rocketmq.StdProducerConfig("example").Build()

		err = producerClient.Start()
		Expect(err).Should(BeNil())

		for i := 0; i < 10; i++ {
			msg := "d" + strconv.Itoa(i)
			err = producerClient.Send([]byte(msg))
			Expect(err).Should(BeNil())
		}

		consumerClient.Close()
		producerClient.Close()
	})
})
