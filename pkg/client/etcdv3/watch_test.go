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

package etcdv3

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	testconf := `
[jupiter.etcdv3.default]
	endpoints = ["http://localhost:2379"]
	timeout = "3s"
	`

	err := conf.LoadFromReader(bytes.NewBufferString(testconf), toml.Unmarshal)
	if err != nil {
		panic(err)
	}
}

func TestClient_WatchPrefix(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test watch",
			fields: fields{
				Client: StdConfig("default").MustSingleton().Client,
				config: StdConfig("default"),
			},
			args: args{
				ctx:    context.Background(),
				prefix: "/test/watch",
			},
			want:    []byte("/test/watch/1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			got, err := client.WatchPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.WatchPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			done := make(chan struct{})
			go func() {
				client.Put(context.Background(), "/test/watch/1", "1")
				time.Sleep(time.Second)
				done <- struct{}{}
			}()

			timeout := time.After(time.Second * 5)

			for {
				select {
				case <-timeout:
					t.Fail()
					return
				case <-done:
					return
				case event := <-got.C():
					assert.Equal(t, event.Kv.Key, tt.want)
				}
			}
		})
	}
}
