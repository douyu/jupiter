package xfreecache

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

var helloReply = &helloworldv1.SayHiResponse{
	Error: 0,
	Msg:   "success",
	Data: &helloworldv1.SayHiResponse_Data{
		Name:      "testName",
		AgeNumber: 18,
	},
}

/*
encoding/json
*/
func BenchmarkDecodeStdStructMedium(b *testing.B) {
	res, _ := json.Marshal(helloReply)
	b.ReportAllocs()
	var data helloworldv1.SayHiResponse
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(res, &data)
	}
}

func BenchmarkEncodeStdStructMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(helloReply)
	}
}

func BenchmarkDecodeJsoniterStructMedium(b *testing.B) {
	res, _ := jsoniter.Marshal(helloReply)
	b.ReportAllocs()
	var data helloworldv1.SayHiResponse
	for i := 0; i < b.N; i++ {
		_ = jsoniter.Unmarshal(res, &data)
	}
}

func BenchmarkEncodeJsoniterStructMedium(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = jsoniter.Marshal(helloReply)
	}
}

func BenchmarkDecodeProto(b *testing.B) {
	res, _ := proto.Marshal(helloReply)
	b.ReportAllocs()
	var data helloworldv1.SayHiResponse
	for i := 0; i < b.N; i++ {
		_ = proto.Unmarshal(res, &data)
	}
}

func BenchmarkEncodeProto(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = proto.Marshal(helloReply)
	}
}

func BenchmarkDecodeProtoWithReflect(b *testing.B) {
	res, _ := proto.Marshal(helloReply)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = unmarshal[*helloworldv1.SayHiResponse](res)
	}
}

func BenchmarkEncodeProtoWithReflect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = marshal[*helloworldv1.SayHiResponse](helloReply)
	}
}

func BenchmarkDecodeProtoWithReflectAndPool(b *testing.B) {
	pool := getPool[*helloworldv1.SayHiResponse]()
	res, _ := proto.Marshal(helloReply)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = unmarshalWithPool[*helloworldv1.SayHiResponse](res, pool)
	}
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
				msgN := reflect.New(msgType).Interface().(proto.Message)
				return msgN
			},
		}
		pools.Store(msgType.String(), pool)
		return pool
	}
	return nil
}

// 反序列化，如果是pb格式，则使用proto序列化 使用sync.Pool-存在并发问题
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
