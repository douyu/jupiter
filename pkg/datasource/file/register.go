package file

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/datasource/manager"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

func init() {
	manager.Register(DataSourceFile, func() conf.DataSource {
		var (
			watchConfig = flag.Bool("watch")
			configAddr  = flag.String("config")
		)
		if configAddr == "" {
			xlog.Panic("new file dataSource, configAddr is empty")
			return nil
		}
		return NewDataSource(configAddr, watchConfig)
	})
	manager.DefaultScheme = DataSourceFile
}
