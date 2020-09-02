// Copyright 2020 Douyu
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

package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/client/rocketmq"
	"github.com/douyu/jupiter/pkg/xlog"
)

// run: go run main.go -config=config.toml
type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.exampleRocketMQProducer,
		eng.exampleRocketMQConsumer,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func main() {
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (eng *Engine) exampleRocketMQConsumer() (err error) {
	consumerClient := rocketmq.StdPushConsumerConfig("configName").Build()
	defer func() {
		if consumerClient.Enable {
			_ = consumerClient.Close()
		}
	}()
	consumerClient.Subscribe(consumerClient.ConsumerConfig.Topic, func(ctx context.Context, ext *primitive.MessageExt) error {
		fmt.Println("msg...", string(ext.Message.Body))
		fmt.Println("msg topic...", string(ext.Message.Topic))
		fmt.Println("msg topic tag...", string(ext.Message.GetTags()))
		return nil
	})
	err = consumerClient.Start()
	return
}

func (eng *Engine) exampleRocketMQProducer() (err error) {
	producerClient := rocketmq.StdProducerConfig("configName").Build()
	defer func() {
		_ = producerClient.Close()
	}()

	err = producerClient.Start()
	if err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		msg := "a" + strconv.Itoa(i)
		err = producerClient.Send([]byte(msg))
	}
	return
}
