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

package rocketmq

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/application"
	"github.com/douyu/jupiter/pkg/governor"
	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var _producers = &sync.Map{}
var _consumers = &sync.Map{}

func init() {
	primitive.PanicHandler = func(i interface{}) {
		stack := make([]byte, 1024)
		length := runtime.Stack(stack, true)
		fmt.Fprint(os.Stderr, "[rocketmq panic recovery]\n", string(stack[:length]))
		xlog.Jupiter().Error("rocketmq panic recovery", zap.Any("error", i))
	}

	governor.HandleFunc("/debug/rocketmq/stats", func(w http.ResponseWriter, r *http.Request) {
		type rocketmqStatus struct {
			application.RuntimeStats
			RocketMQs map[string]interface{} `json:"rocketmqs"`
			FlowInfo  map[string]FlowInfo    `json:"flowInfo"`
		}

		var rets = rocketmqStatus{
			RuntimeStats: application.NewRuntimeStats(),
			RocketMQs:    make(map[string]interface{}),
			FlowInfo:     make(map[string]FlowInfo),
		}

		_producers.Range(func(key interface{}, val interface{}) bool {
			name := key.(string)
			cc := val.(*Producer)
			rets.RocketMQs[name] = map[string]interface{}{
				"role":   "producer",
				"config": cc.ProducerConfig,
			}
			rets.FlowInfo[fmt.Sprintf("%s_%s", name, cc.fInfo.GroupType)] = cc.fInfo
			return true
		})

		_consumers.Range(func(key interface{}, val interface{}) bool {
			name := key.(string)
			cc := val.(*PushConsumer)
			rets.RocketMQs[name] = map[string]interface{}{
				"config": cc.ConsumerConfig,
				"role":   "consumer",
			}
			rets.FlowInfo[name+"_"+cc.fInfo.GroupType] = cc.fInfo
			return true
		})

		_ = jsoniter.NewEncoder(w).Encode(rets)
	})
}

func GetProducer(name string) *Producer {
	if ins, ok := _producers.Load(name); ok {
		return ins.(*Producer)
	}
	return nil
}

// Get ...
func GetConsumer(name string) *PushConsumer {
	if ins, ok := _consumers.Load(name); ok {
		return ins.(*PushConsumer)
	}

	return nil
}

// Invoker ...
func InvokerProducer(name string) *Producer {
	if client := GetProducer(name); client != nil {
		return client
	}

	return StdNewProducer(name)
}
