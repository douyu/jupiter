package mongox

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
)

func TestMain(m *testing.M) {
	conf.LoadFromReader(bytes.NewBufferString(`
	[jupiter.mongo.demo]
		dsn = ""
		socketTimeout = "5s"
		poolLimit = 100
	`), toml.Unmarshal)
}

func TestStdConfig(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		// TODO: Add test cases.
		{
			name: "std config",
			args: args{
				name: "demo",
			},
			want: Config{
				DSN:           "",
				SocketTimeout: time.Second * 5,
				PoolLimit:     100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StdConfig(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StdConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawConfig(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		// TODO: Add test cases.
		{
			name: "raw config",
			args: args{
				key: "jupiter.mongo.demo",
			},
			want: Config{
				DSN:           "",
				SocketTimeout: time.Second * 5,
				PoolLimit:     100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RawConfig(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
