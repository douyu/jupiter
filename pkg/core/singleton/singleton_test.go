package singleton

import (
	"reflect"
	"testing"

	"github.com/douyu/jupiter/pkg/core/constant"
)

func TestStore(t *testing.T) {
	{
		type args struct {
			module constant.Module
			key    string
			val    interface{}
		}
		tests := []struct {
			name string
			args args
		}{
			// TODO: Add test cases.
			{
				name: "test1",
				args: args{
					module: constant.ModuleClientEtcd,
					key:    "key",
					val:    "5test1",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				Store(tt.args.module, tt.args.key, tt.args.val)
			})
		}
	}

	{
		type args struct {
			module constant.Module
			key    string
		}
		tests := []struct {
			name  string
			args  args
			want  interface{}
			want1 bool
		}{
			// TODO: Add test cases.
			{
				name: "test1",
				args: args{
					module: constant.ModuleClientEtcd,
					key:    "key",
				},
				want:  "5test1",
				want1: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, got1 := Load(tt.args.module, tt.args.key)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Load() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("Load() got1 = %v, want %v", got1, tt.want1)
				}
			})
		}
	}
}

func Test_genkey(t *testing.T) {
	type args struct {
		module constant.Module
		key    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				module: constant.ModuleClientEtcd,
				key:    "key",
			},
			want: "5key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genkey(tt.args.module, tt.args.key); got != tt.want {
				t.Errorf("genkey() = %v, want %v", got, tt.want)
			}
		})
	}
}
