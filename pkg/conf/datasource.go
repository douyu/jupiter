package conf

import (
	"errors"
	"io"
	"log"
	"net/url"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/flag"
)

var (
	//ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	datasourceBuilders   map[string]DataSourceCreatorFunc
	configDecoder        map[string]Unmarshaller
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() DataSource

// DataSource ...
type DataSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

func init() {
	datasourceBuilders = make(map[string]DataSourceCreatorFunc)
	flag.Register(&flag.StringFlag{Name: "config", Usage: "--config=config.toml", Action: func(key string, fs *flag.FlagSet) {
		var configAddr = fs.String(key)
		log.Printf("read config: %s", configAddr)
		datasource, err := NewDataSource(configAddr)
		if err != nil {
			log.Fatalf("build datasource[%s] failed: %v", configAddr, err)
		}
		if err := LoadFromDataSource(datasource, toml.Unmarshal); err != nil {
			log.Fatalf("load config from datasource[%s] failed: %v", configAddr, err)
		}
		log.Printf("load config from datasource[%s] completely!", configAddr)
	}})
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator DataSourceCreatorFunc) {
	datasourceBuilders[scheme] = creator
}

// CreateDataSource creates a dataSource witch has been registered
// func CreateDataSource(scheme string) (conf.DataSource, error) {
// 	creatorFunc, exist := registry[scheme]
// 	if !exist {
// 		return nil, ErrInvalidDataSource
// 	}
// 	return creatorFunc(), nil
// }

//NewDataSource ..
func NewDataSource(configAddr string) (DataSource, error) {
	if configAddr == "" {
		return nil, ErrConfigAddr
	}
	urlObj, err := url.Parse(configAddr)
	if err != nil {
		return nil, err
	}

	var scheme = urlObj.Scheme
	if scheme == "" {
		scheme = "file"
	}

	creatorFunc, exist := datasourceBuilders[scheme]
	if !exist {
		return nil, ErrInvalidDataSource
	}
	return creatorFunc(), nil
}
