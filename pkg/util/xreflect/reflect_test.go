package xreflect_test

import (
	"testing"

	"github.com/douyu/jupiter/pkg/util/xreflect"
)

func TestIn(t *testing.T) {
	type args struct {
		value     interface{}
		container interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key in map",
			args: args{value: "foo", container: map[string]interface{}{"foo": "bar"}},
			want: true,
		},
		{
			name: "key not in map",
			args: args{value: "foo", container: map[string]interface{}{"foo2": "bar"}},
			want: false,
		},
		{
			name: "value in slice",
			args: args{value: "foo", container: []string{"foo", "bar"}},
			want: true,
		},
		{
			name: "value not in slice",
			args: args{value: "foo", container: []string{"foo2", "bar"}},
			want: false,
		},
		{
			name: "value in array",
			args: args{value: "foo", container: [...]string{"foo", "bar"}},
			want: true,
		},
		{
			name: "value not in array",
			args: args{value: "foo", container: [...]string{"foo2", "bar"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xreflect.In(tt.args.value, tt.args.container); got != tt.want {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverride(t *testing.T) {
	type User struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	var left = User{Name: "foo"}
	var right = User{Address: "unknown"}
	type args struct {
		left  interface{}
		right interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "merge map string",
			args: args{left: &left, right: &right},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xreflect.Override(tt.args.left, tt.args.right); got != tt.want {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
			t.Logf("left: %v, right: %v", left, right)
		})
	}
}
