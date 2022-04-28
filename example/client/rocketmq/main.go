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
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"strconv"
	"time"

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
	consumerClient.RegisterBatchMessage(func(ctx context.Context, msgs ...*primitive.MessageExt) error {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}
		time.Sleep(50 * time.Millisecond)
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
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		msg := "a" + strconv.Itoa(i)
		producerClient.SendWithContext(ctx, []byte(msg))
	}

	return
}
