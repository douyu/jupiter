package xfreecache

import (
	"encoding/json"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"reflect"
	"testing"

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
		_ = protoUnmarshal[*helloworldv1.SayHiResponse](res)
	}
}

func protoUnmarshal[T any](body []byte) T {
	var value T

	if msg, ok := any(value).(proto.Message); ok { // Constrained to proto.Message
		// Peek the type inside T (as T= *SomeProtoMsgType)
		msgType := reflect.TypeOf(msg).Elem()

		// Make a new one, and throw it back into T
		msg = reflect.New(msgType).Interface().(proto.Message)

		_ = proto.Unmarshal(body, msg)
		value = msg.(T)
	}

	return value
}
