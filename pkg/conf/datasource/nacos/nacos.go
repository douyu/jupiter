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
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type nacosDataSource struct {
	client config_client.IConfigClient
	group  string
	dataID string

	changed chan struct{}
}

// NewDataSource creates an nacos DataSource
func NewDataSource(client config_client.IConfigClient, group, dataID string, watch bool) conf.DataSource {
	ds := &nacosDataSource{
		client:  client,
		group:   group,
		dataID:  dataID,
		changed: make(chan struct{}, 1),
	}
	if watch {
		_ = ds.client.ListenConfig(vo.ConfigParam{
			Group:  ds.group,
			DataId: ds.dataID,
			OnChange: func(namespace, group, dataId, data string) {
				ds.changed <- struct{}{}
			},
		})
	}
	return ds
}

// ReadConfig reads config content from nacos
func (ds *nacosDataSource) ReadConfig() ([]byte, error) {
	configData, err := ds.client.GetConfig(vo.ConfigParam{
		Group:  ds.group,
		DataId: ds.dataID,
	})
	if err != nil {
		return nil, err
	}

	return []byte(configData), nil
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ds *nacosDataSource) IsConfigChanged() <-chan struct{} {
	return ds.changed
}

// Close stops watching the config changed
func (ds *nacosDataSource) Close() error {
	_ = ds.client.CancelListenConfig(vo.ConfigParam{
		Group:  ds.group,
		DataId: ds.dataID,
	})
	ds.client.CloseClient()
	close(ds.changed)
	return nil
}
