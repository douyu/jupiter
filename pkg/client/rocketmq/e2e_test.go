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
	. "github.com/onsi/ginkgo"
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

		defer func() {
			if consumerClient.Enable {
				consumerClient.Close()
			}
		}()

		count := int32(0)

		consumerClient.Subscribe(consumerClient.ConsumerConfig.Topic, func(ctx context.Context, ext *primitive.MessageExt) error {
			atomic.AddInt32(&count, 1)
			fmt.Println("msg...", string(ext.Message.Body), string(ext.Message.Topic), string(ext.Message.GetTags()), atomic.LoadInt32(&count))

			return nil
		})
		err := consumerClient.Start()
		Expect(err).Should(BeNil())

		producerClient := rocketmq.StdProducerConfig("example").Build()
		defer producerClient.Close()

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
	})
})
