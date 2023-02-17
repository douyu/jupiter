package xfreecache

import (
	"encoding/json"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

// 序列化，如果是pb格式，则使用proto序列化
func marshal[T any](cacheData T) (data []byte, err error) {
	if msg, ok := any(cacheData).(proto.Message); ok {
		data, err = proto.Marshal(msg)
	} else {
		data, err = json.Marshal(cacheData)
	}
	return
}

var pools sync.Map

func getPool[T any]() *sync.Pool {
	var value T
	if msg, ok := any(value).(proto.Message); ok {
		msgType := reflect.TypeOf(msg).Elem()
		if pool, ok2 := pools.Load(msgType.String()); ok2 {
			return pool.(*sync.Pool)
		}

		pool := &sync.Pool{
			New: func() any {
				// Make a new one, and throw it back into T
				msg = reflect.New(msgType).Interface().(proto.Message)
				return msg
			},
		}
		pools.Store(msgType.String(), pool)
		return pool
	}
	return nil
}

// 反序列化，如果是pb格式，则使用proto序列化 使用sync.Pool
func unmarshalWithPool[T any](body []byte, pool *sync.Pool) (value T, err error) {
	if _, ok := any(value).(proto.Message); ok { // Constrained to proto.Message
		msg := pool.Get().(proto.Message)
		err = proto.Unmarshal(body, msg)
		value = msg.(T)
		pool.Put(msg)
	} else {
		err = json.Unmarshal(body, &value)
	}
	return
}
