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
	"testing"
)

func TestConfigBuild(t *testing.T) {
	// TODO(gorexlv): add etcd ci
	// config := &Config{
	// 	Endpoints:      []string{"127.0.0.1:2379"},
	// 	ConnectTimeout: time.Second * 3,
	// 	Secure:         false,
	// 	logger:         xlog.DefaultConfig().Build(),
	// }
	// etcdClient := config.Build()
	// _, err := etcdClient.Put(context.Background(), "etcdTest", "test")
	// assert.Equal(t, nil, err)
	// fmt.Println("etcdClient", etcdClient)

	// // test get
	// kv, err := etcdClient.GetKeyValue(context.Background(), "etcdTest")
	// assert.Equal(t, nil, err)
	// assert.Equal(t, "etcdTest", string(kv.Key))
	// assert.Equal(t, "test", string(kv.Value))

	// // test getprefix
	// kvMap, err := etcdClient.GetPrefix(context.Background(), "etcd")
	// assert.Equal(t, nil, err)
	// assert.Equal(t, map[string]string{"etcdTest": "test"}, kvMap)

	// // test del
	// _, err = etcdClient.Delete(context.Background(), "etcdTest")
	// assert.Equal(t, nil, err)
	// kv, err = etcdClient.GetKeyValue(context.Background(), "etcdTest")
	// assert.Equal(t, nil, err)
	// assert.Equal(t, (*mvccpb.KeyValue)(nil), kv)
}
