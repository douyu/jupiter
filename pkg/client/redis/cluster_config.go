package redis

import (
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

// ClusterOptions 集群配置信息
type ClusterOptions struct {
	Config
}

// ClusterConfig 集群模式
func ClusterConfig(name string) *ClusterOptions {
	return RawClusterConfig(constant.ConfigKey("redis", name, "cluster"))
}

func RawClusterConfig(key string) *ClusterOptions {
	option := DefaultClusterOption()
	if err := cfg.UnmarshalKey(key, &option, cfg.TagName("toml")); err != nil {
		option.logger.Panic("unmarshal config:"+key, xlog.FieldErr(err), xlog.FieldName(key), xlog.FieldExtMessage(option))
	}

	if len(option.Addr) == 0 {
		option.logger.Panic("no cluster addr set:"+key, xlog.FieldName(key), xlog.FieldExtMessage(option))
	}
	option.name = key
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint(key, option)
	}

	return option
}

// DefaultClusterOption default option ...
func DefaultClusterOption() *ClusterOptions {
	return &ClusterOptions{Config: *DefaultConfig()}
}
