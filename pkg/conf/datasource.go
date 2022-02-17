package conf

import (
	"errors"
	"io"
	"net/url"
)

var (
	//ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	datasourceBuilders   = make(map[string]DataSourceCreatorFunc)
	// configDecoder        = make(map[string]Unmarshaller)
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() DataSource

// DataSource ...
type DataSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
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
