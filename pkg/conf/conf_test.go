// Copyright 2021 douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestConf(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		LoadFromReader(bytes.NewBuffer([]byte(`
		[server]
		[server.http]
		[server.http.addr]
			port = 8080
			addr = "localhost"
		`)), toml.Unmarshal)

		if GetInt("server.http.addr.port") != 8080 {
			t.Fatal("get int failed")
		}

		if GetString("server.http.addr.addr") != "localhost" {
			t.Fatal("get string failed")
		}

		type Addr struct {
			Port int    `toml:"port"`
			Addr string `toml:"addr"`
		}

		addr := Addr{}
		err := UnmarshalKey("server.http.addr", &addr)
		if err != nil {
			t.Fatal(err)
		}

		if addr.Port != 8080 {
			t.Fatal("unmarshal failed")
		}

		if addr.Addr != "localhost" {
			t.Fatal("unmarshal failed")
		}

		type Cfg struct {
			Server struct {
				HTTP struct {
					Addr Addr
				}
			}
		}

		cfg := Cfg{}
		err = UnmarshalKey("", &cfg)
		if err != nil {
			t.Fatal(err)
		}

		if cfg.Server.HTTP.Addr.Port != 8080 {
			t.Fatal("unmarshal failed")
		}

		if cfg.Server.HTTP.Addr.Addr != "localhost" {
			t.Fatal("unmarshal failed")
		}
	})
}

func TestGetInt64Slice(t *testing.T) {

	LoadFromReader(bytes.NewBufferString(`
	[test]
		ids = [1000]
	`), toml.Unmarshal)
	defer Reset()

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "case 1",
			args: args{
				key: "test.ids",
			},
			want: []int64{1000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInt64Slice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInt64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}
