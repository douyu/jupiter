package xhttprule

import (
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/douyu/jupiter/proto/testproto/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestGetHTTPRules(t *testing.T) {

	type args struct {
		m protoreflect.MethodDescriptor
	}
	tests := []struct {
		name string
		args args
		want []*HTTPRule
	}{
		{
			name: "case 1: normal post",
			args: args{
				m: testproto.File_testproto_v1_hello_proto.Services().Get(0).Methods().Get(0),
			},
			want: []*HTTPRule{
				{
					Method: "POST",
					URI:    "/v1/helloworld.Greeter/SayHello",
					Body:   "*",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetHTTPRules(tt.args.m))
		})
	}
}

func TestHTTPRule_GetPlainURI(t *testing.T) {
	type fields struct {
		Method       string
		URI          string
		Body         string
		ResponseBody string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "case 1",
			fields: fields{
				Method: "POST",
				URI:    "/v1/helloworld.Greeter/SayHello",
				Body:   "*",
			},
			want: "/v1/helloworld.Greeter/SayHello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPRule{
				Method:       tt.fields.Method,
				URI:          tt.fields.URI,
				Body:         tt.fields.Body,
				ResponseBody: tt.fields.ResponseBody,
			}
			if got := h.GetPlainURI(); got != tt.want {
				t.Errorf("HTTPRule.GetPlainURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRule_GetVariables(t *testing.T) {
	type fields struct {
		Method       string
		URI          string
		Body         string
		ResponseBody string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "case 1",
			fields: fields{
				Method: "POST",
				URI:    "/v1/helloworld.Greeter/SayHello",
				Body:   "*",
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPRule{
				Method:       tt.fields.Method,
				URI:          tt.fields.URI,
				Body:         tt.fields.Body,
				ResponseBody: tt.fields.ResponseBody,
			}
			if got := h.GetVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPRule.GetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}
