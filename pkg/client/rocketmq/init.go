package rocketmq

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/douyu/jupiter/pkg/application"
	jsoniter "github.com/json-iterator/go"
)

var _producers = &sync.Map{}
var _consumers = &sync.Map{}

func init() {
	http.HandleFunc("/debug/rocketmq/stats", func(w http.ResponseWriter, r *http.Request) {
		type rocketmqStatus struct {
			application.RuntimeStats
			RocketMQs map[string]interface{} `json:"rocketmqs"`
			FlowInfo  map[string]FlowInfo    `json:"flowInfo"`
		}

		var rets = rocketmqStatus{
			RuntimeStats: application.NewRuntimeStats(),
			RocketMQs:    make(map[string]interface{}, 0),
			FlowInfo:     make(map[string]FlowInfo, 0),
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
