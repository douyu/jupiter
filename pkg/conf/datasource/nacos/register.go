package nacos

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/util/xnet"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	xcast "github.com/spf13/cast"
)

// DataSourceNacos defines nacos scheme
const DataSourceNacos = "nacos"

func init() {
	conf.Register(DataSourceNacos, func() conf.DataSource {
		var (
			configAddr = flag.String("config")
			watch      = flag.Bool("watch")
		)
		if configAddr == "" {
			xlog.Jupiter().Panic("new nacos dataSource, configAddr is empty")
			return nil
		}
		// configAddr is a string in this format:
		// nacos://ip:port?dataId=xx&group=xx&namespaceId=xx&timeout=10000&accessKey=xx&secretKey=xx&notLoadCacheAtStart=true&updateCacheWhenEmpty=true
		urlObj, err := xnet.ParseURL(configAddr)
		if err != nil {
			xlog.Jupiter().Panic("parse configAddr error", xlog.Any("error", err))
			return nil
		}
		// create clientConfig
		clientConfig := constant.ClientConfig{
			TimeoutMs:            urlObj.QueryUint64("timeout", 10000),
			NotLoadCacheAtStart:  urlObj.QueryBool("notLoadCacheAtStart", true),
			UpdateCacheWhenEmpty: urlObj.QueryBool("updateCacheWhenEmpty", true),
			NamespaceId:          urlObj.Query().Get("namespaceId"),
			AccessKey:            urlObj.Query().Get("accessKey"),
			SecretKey:            urlObj.Query().Get("secretKey"),
		}
		// create serverConfigs
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr: urlObj.HostName,
				Port:   getPort(urlObj.Port, 8848),
			},
		}
		// create config client
		client, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			xlog.Jupiter().Panic("create config client error", xlog.Any("error", err))
			return nil
		}
		return NewDataSource(client, urlObj.Query().Get("group"), urlObj.Query().Get("dataId"), watch)
	})
}

func getPort(port string, expect uint64) uint64 {
	ret, err := xcast.ToUint64E(port)
	if err != nil {
		return expect
	}
	return ret
}
