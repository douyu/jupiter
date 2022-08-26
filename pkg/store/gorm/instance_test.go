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
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Birthday time.Time
	Age      int
	Name     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
}

func (u User) TableName() string {
	return "users"
}

func TestStdNew(t *testing.T) {
	type args struct {
		name string
		opts []interface{}
	}
	tests := []struct {
		name    string
		args    args
		config  string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "std new",
			args: args{
				name: "demo",
				opts: []interface{}{},
			},
			wantErr: false,
			config: `
			[jupiter.mysql.demo]
				dsn="root:123456@tcp(localhost:3306)/mysql?timeout=20s&readTimeout=20s"
				debug=true
				maxIdleConns=50
				connMaxLifeTime="20m"
				level="panic"
				slowThreshold="400ms"
				dialTimeout="1s"
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			defer cfg.Reset()
			if tt.wantErr {
				assert.Panics(t, func() { StdConfig(tt.args.name).MustBuild() })
			} else {
				assert.NotPanics(t, func() { StdConfig(tt.args.name).MustBuild() })
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			defer cfg.Reset()
			assert.NotPanics(t, func() { StdConfig(tt.args.name).MustBuild() })
		})
	}
}
