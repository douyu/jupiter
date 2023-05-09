package xfreecache

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"reflect"
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

// 反序列化，如果是pb格式，则使用proto序列化
func unmarshal[T any](body []byte) (value T, err error) {
	if msg, ok := any(value).(proto.Message); ok { // Constrained to proto.Message
		// Peek the type inside T (as T= *SomeProtoMsgType)
		msgType := reflect.TypeOf(msg).Elem()

		// Make a new one, and throw it back into T
		msg = reflect.New(msgType).Interface().(proto.Message)

		err = proto.Unmarshal(body, msg)
		value = msg.(T)
	} else {
		err = json.Unmarshal(body, &value)
	}
	return
}
