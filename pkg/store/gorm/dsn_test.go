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

package gorm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDSN(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantCfg *DSN
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				dsn: "user:password@tcp(127.0.0.1:9090)/config?aa=aa&&bb=bb&&cc=cc",
			},
			wantCfg: &DSN{
				User:     "user",
				Password: "password",
				Net:      "tcp",
				Addr:     "127.0.0.1:9090",
				DBName:   "config",
				Params:   map[string]string{"aa": "aa", "bb": "bb", "cc": "cc"},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				dsn: ":@tcp(127.0.0.1)/config?aa=aa&&bb=bb&&cc=cc",
			},
			wantCfg: &DSN{
				User:     "",
				Password: "",
				Net:      "tcp",
				Addr:     "127.0.0.1",
				DBName:   "config",
				Params:   map[string]string{"aa": "aa", "bb": "bb", "cc": "cc"},
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				dsn: "tcp(127.0.0.1)/config?aa=aa&&bb=bb&&cc=cc",
			},
			wantCfg: &DSN{
				User:     "",
				Password: "",
				Net:      "tcp",
				Addr:     "127.0.0.1",
				DBName:   "config",
				Params:   map[string]string{"aa": "aa", "bb": "bb", "cc": "cc"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := ParseDSN(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("ParseDSN() = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}

func Test_parseDSNParams(t *testing.T) {
	type args struct {
		cfg    *DSN
		params string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				cfg: &DSN{
					User:     "user",
					Password: "password",
					Net:      "tcp",
					Addr:     "127.0.0.1:9090",
					DBName:   "config",
					Params:   nil,
				},
				params: "aa=aa&bb=bb&cc=cc",
			},
			wantErr: false,
			want:    map[string]string{"aa": "aa", "bb": "bb", "cc": "cc"},
		},
		{
			name: "2",
			args: args{
				cfg: &DSN{
					User:     "user",
					Password: "password",
					Net:      "tcp",
					Addr:     "127.0.0.1:9090",
					DBName:   "config",
					Params:   nil,
				},
				params: "aa=[a1,a2]&bb=bb|bb1&cc=cc,cc1",
			},
			wantErr: false,
			want:    map[string]string{"aa": "[a1,a2]", "bb": "bb|bb1", "cc": "cc,cc1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseDSNParams(tt.args.cfg, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("parseDSNParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.args.cfg.Params, tt.want)
		})
	}
}
