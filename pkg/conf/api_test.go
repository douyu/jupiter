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

package conf

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadFromReader(t *testing.T) {
	type args struct {
		r            io.Reader
		unmarshaller Unmarshaller
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "toml",
			args: args{
				r: bytes.NewBufferString(`
				[server]
				[server.http]
				[server.http.addr]
					port = 8080
					addr = "localhost"
				`),
				unmarshaller: toml.Unmarshal,
			},
			wantErr: false,
		},
		{
			name: "json",
			args: args{
				r: bytes.NewBufferString(`
				{
					"server": {
						"http": {
							"addr": {
								"port": 8080,
								"addr": "localhost"
							}
						}
					}
				}
				`),
				unmarshaller: json.Unmarshal,
			},
			wantErr: false,
		},
		{
			name: "yml",
			args: args{
				r: bytes.NewBufferString(`
server:
  http:
    addr:
      port: 8080
      addr: "localhost"
`),
				unmarshaller: yaml.Unmarshal,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadFromReader(tt.args.r, tt.args.unmarshaller); (err != nil) != tt.wantErr {
				t.Errorf("LoadFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, GetInt("server.http.addr.port"), 8080)
			assert.Equal(t, GetString("server.http.addr.addr"), "localhost")
			Reset()
			assert.Equal(t, GetInt("server.http.addr.port"), 0)
			assert.Equal(t, GetString("server.http.addr.addr"), "")
		})
	}
}
