package juno

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Defines juno scheme
const (
	DataSourceJuno = "juno"
)

func init() {
	dataSourceCreator := func() conf.DataSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			xlog.Jupiter().Panic("new http dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	}
	conf.Register(DataSourceJuno, dataSourceCreator)
}
