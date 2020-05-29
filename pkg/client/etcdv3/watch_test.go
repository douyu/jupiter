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
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestEtcdWatch(t *testing.T) {
	config := &Config{
		Endpoints:      []string{"127.0.0.1:2379"},
		ConnectTimeout: time.Second * 3,
		Secure:         false,
		logger:         xlog.DefaultConfig().Build(),
	}
	etcdClient := config.Build()

	watcher, err := etcdClient.NewWatch("etcdTest")
	assert.Equal(t, nil, err)

	go func() {
		ticKer := time.NewTicker(time.Second * 5)
		for {
			select {
			case event := <-watcher.C():
				t.Log("event", event)
			case <-ticKer.C:
				t.Log("ticker")

			}
		}
	}()

	cases := []struct {
		Type string `json:"type"`
		Key  string `json:"key"`
		Val  string `json:"val"`
	}{
		{
			Type: "put",
			Key:  "etcdTest",
			Val:  "test",
		},
		{
			Type: "put",
			Key:  "etcdTest",
			Val:  "test1",
		},
		{
			Type: "del",
			Key:  "etcdTest",
		},
	}

	for _, info := range cases {
		switch info.Type {
		case "put":
			_, err = etcdClient.Put(context.Background(), info.Key, info.Val)
			assert.Equal(t, nil, err)
		case "del":
			_, err = etcdClient.Delete(context.Background(), info.Key)
			assert.Equal(t, nil, err)
		}
	}

	time.Sleep(time.Second * 3)
}
