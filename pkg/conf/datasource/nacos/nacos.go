package nacos

import (
	"log"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xgo"
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
		client: client,
		group:  group,
		dataID: dataID,
	}
	if watch {
		ds.changed = make(chan struct{}, 1)
		xgo.Go(ds.watch)
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

func (ds *nacosDataSource) watch() {
	ds.client.ListenConfig(vo.ConfigParam{
		Group:  ds.group,
		DataId: ds.dataID,
		OnChange: func(namespace, group, dataId, data string) {
			log.Println("nacos config changed: ", data)
			ds.changed <- struct{}{}
		},
	})
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ds *nacosDataSource) IsConfigChanged() <-chan struct{} {
	return ds.changed
}

// Close stops watching the config changed
func (ds *nacosDataSource) Close() error {
	ds.client.CancelListenConfig(vo.ConfigParam{
		Group:  ds.group,
		DataId: ds.dataID,
	})
	ds.client.CloseClient()
	close(ds.changed)
	return nil
}
