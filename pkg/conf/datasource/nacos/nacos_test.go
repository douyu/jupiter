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

package nacos

import (
	"sync"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/conf/datasource/nacos/mock"
	"github.com/golang/mock/gomock"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
)

var (
	wg         sync.WaitGroup
	localParam = vo.ConfigParam{
		DataId:  "data-id",
		Group:   "group",
		Content: "hello world",
	}
	newContent = "hello-world-changed"
)

func TestReadConfig(t *testing.T) {

	t.Run("with watch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewMockIConfigClient(ctrl)

		client.EXPECT().CancelListenConfig(gomock.Any()).Return(nil)
		client.EXPECT().CloseClient().Return()
		client.EXPECT().GetConfig(gomock.Any()).Return(localParam.Content, nil)
		client.EXPECT().ListenConfig(gomock.Any()).DoAndReturn(func(param vo.ConfigParam) error {
			go func() {
				time.Sleep(time.Second)
				client.EXPECT().GetConfig(gomock.Any()).Return(newContent, nil)
				param.OnChange("namespace", localParam.Group, localParam.DataId, newContent)
			}()

			return nil
		})

		ds := NewDataSource(client, localParam.Group, localParam.DataId, true)

		content, err := ds.ReadConfig()
		assert.Nil(t, err)
		assert.Equal(t, localParam.Content, string(content))
		t.Logf("read config: %s", content)

		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ds.IsConfigChanged() {
				content, err := ds.ReadConfig()
				assert.Nil(t, err)
				assert.Equal(t, newContent, string(content))
				t.Logf("read new config: %s", content)
				ds.Close()
			}
		}()

		wg.Wait()
	})

	t.Run("without with", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewMockIConfigClient(ctrl)
		client.EXPECT().CancelListenConfig(gomock.Any()).Return(nil)
		client.EXPECT().CloseClient().Return()
		client.EXPECT().GetConfig(gomock.Any()).Return(localParam.Content, nil)
		ds := NewDataSource(client, localParam.Group, localParam.DataId, false)
		defer ds.Close()

		content, err := ds.ReadConfig()
		assert.Nil(t, err)
		assert.Equal(t, localParam.Content, string(content))
		t.Logf("read config: %s", content)
	})
}
