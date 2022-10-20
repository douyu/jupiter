// Copyright 2022 Douyu
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

package gorm

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/stretchr/testify/assert"
)

func TestStdConfig(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		want   *Config
		config string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				name: "demo",
			},
			want: &Config{
				Name:            constant.ConfigKey("mysql.demo"),
				DSN:             "this is dsn",
				Debug:           true,
				MaxIdleConns:    100,
				MaxOpenConns:    100,
				ConnMaxLifetime: time.Second,
				OnDialError:     "panic",
				SlowThreshold:   time.Millisecond * 100,
				DialTimeout:     3 * time.Second,

				raw:           nil,
				Retry:         2,
				RetryWaitTime: time.Millisecond * 200,
			},
			config: `
				[jupiter.mysql.demo]
					dsn="this is dsn"
					debug=true
					maxIdleConns=100
					maxOpenConns=100
					connMaxLifetime="1s"
					level="panic"
					slowThreshold="100ms"
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
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
		name   string
		args   args
		want   *Config
		config string
		result bool
	}{
		// TODO: Add test cases.
		{
			name: "demo",
			args: args{key: constant.ConfigKey("mysql.demo1")},
			want: &Config{
				Name:            constant.ConfigKey("mysql.demo1"),
				DSN:             "this is dsn",
				Debug:           true,
				MaxIdleConns:    100,
				MaxOpenConns:    100,
				ConnMaxLifetime: time.Second,
				OnDialError:     "panic",
				SlowThreshold:   time.Millisecond * 100,
				DialTimeout:     3 * time.Second,
				raw:             nil,
				Retry:           3,
				RetryWaitTime:   time.Millisecond * 200,
			},
			config: `
				[jupiter.mysql.demo1]
					dsn="this is dsn"
					debug=true
					maxIdleConns=100
					maxOpenConns=100
					connMaxLifetime="1s"
					level="panic"
					slowThreshold="100ms"
					retry=3
			`,
			result: true,
		},
		{
			name: "default",
			args: args{key: constant.ConfigKey("mysql.demo2")},
			want: &Config{
				Name:            constant.ConfigKey("mysql.demo2"),
				DSN:             "this is dsn",
				Debug:           true,
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: time.Second * 300,
				OnDialError:     "panic",
				SlowThreshold:   time.Millisecond * 500,
				DialTimeout:     3 * time.Second,
				raw:             nil,
				Retry:           2,
				RetryWaitTime:   time.Millisecond * 200,
			},
			config: `
				[jupiter.mysql.demo2]
					dsn="this is dsn"
					debug=true
			`,
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			if got := RawConfig(tt.args.key); !reflect.DeepEqual(got, tt.want) == tt.result {
				t.Errorf("RawConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
		{
			name: "demo",
			want: &Config{
				DSN:             "",
				Debug:           false,
				MaxIdleConns:    10,
				MaxOpenConns:    100,
				ConnMaxLifetime: time.Second * 300,
				OnDialError:     "panic",
				SlowThreshold:   time.Millisecond * 500,
				DialTimeout:     3 * time.Second,
				raw:             nil,
				Retry:           2,
				RetryWaitTime:   time.Millisecond * 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
