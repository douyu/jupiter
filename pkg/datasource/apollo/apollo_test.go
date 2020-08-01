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

package apollo

import (
	"github.com/douyu/jupiter/pkg/datasource/apollo/mockserver"
	"github.com/philchia/agollo/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Println(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Second)
}

func teardown() {
	mockserver.Close()
}

func TestReadConfig(t *testing.T) {
	testData := []string{"value1", "value2"}

	mockserver.Set("application", "key_test", testData[0])
	ds := NewDataSource(&agollo.Conf{
		AppID:          "SampleApp",
		Cluster:        "default",
		NameSpaceNames: []string{"application"},
		MetaAddr:       "localhost:16852",
		CacheDir:       ".",
	}, "application", "key_test")
	value, err := ds.ReadConfig()
	assert.Nil(t, err)
	assert.Equal(t, testData[0], string(value))
	t.Logf("read: %s", value)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mockserver.Set("application", "key_test", testData[1])
		time.Sleep(time.Second * 3)
		ds.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ds.IsConfigChanged() {
			value, err := ds.ReadConfig()
			assert.Nil(t, err)
			assert.Equal(t, testData[1], string(value))
			t.Logf("read: %s", value)
		}
	}()

	wg.Wait()
}
