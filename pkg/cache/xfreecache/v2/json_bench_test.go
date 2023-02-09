package xfreecache

import (
	"encoding/json"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"

	"github.com/douyu/jupiter/proto/testproto"
)

var helloReply = &testproto.HelloReply{
	Message: "benchmarkTest",
	Id64:    -123456,
	Id32:    -12345678,
	Idu64:   123456,
	Idu32:   123456780,
	Name:    []byte("newName"),
	Done:    true,
}

/*
   encoding/json
*/
func BenchmarkDecodeStdStructMedium(b *testing.B) {
	res, _ := json.Marshal(helloReply)
	b.ReportAllocs()
	var data testproto.HelloReply
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
	var data testproto.HelloReply
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
	var data testproto.HelloReply
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
		_ = protoUnmarshal[*testproto.HelloReply](res)
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
