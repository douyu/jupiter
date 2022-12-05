// Copyright 2020 Douyu
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

package tstore

import (
	"bytes"
	"github.com/BurntSushi/toml"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	tablestoreEndponint = os.Getenv("TABLESTORE_ENDPOINT_ENV")
	accessKeyId         = os.Getenv("TABLESTORE_ACCESSKEYID_ENV")
	accessKeySecret     = os.Getenv("TABLESTORE_ACCESSKEYSECRET_ENV")
	tablestoreInstance  = os.Getenv("TABLESTORE_INSTANCE_ENV")
)

func TestStdConfig(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		config string
		want   *Config
	}{
		{
			name: "配置读取",
			args: args{
				name: "demo",
			},
			config: `
			[jupiter.tablestore.demo]
			   debug = false
			   enableAccessLog = false
			   endPoint ="` + tablestoreEndponint + `"
			   instance = "` + tablestoreInstance + `"
			   accessKeyId ="` + accessKeyId + `"
			   accessKeySecret = "` + accessKeySecret + `"
			   requestTimeout = "30s"
			   slowThreshold = "1s"
			   maxIdleConnections = 2000
			`,
			want: &Config{
				Name:               constant.ConfigKey("tablestore.demo"),
				Debug:              false,
				EnableAccessLog:    false,
				EnableMetric:       true,
				EndPoint:           tablestoreEndponint,
				Instance:           tablestoreInstance,
				AccessKeyId:        accessKeyId,
				AccessKeySecret:    accessKeySecret,
				SecurityToken:      "",
				RetryTimes:         1,
				SlowThreshold:      time.Second * 1,
				MaxRetryTime:       time.Second * 5,
				ConnectionTimeout:  time.Second * 15,
				RequestTimeout:     time.Second * 30,
				MaxIdleConnections: 2000,
			},
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
