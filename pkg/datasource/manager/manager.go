package manager

import (
	"errors"
	"github.com/douyu/jupiter/pkg/conf"
)

var (
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	registry             map[string]DataSourceCreatorFunc
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() conf.DataSource

func init() {
	registry = make(map[string]DataSourceCreatorFunc)
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator func() conf.DataSource) {
	registry[scheme] = creator
}

// CreateDataSource creates a dataSource witch has been registered
func CreateDataSource(scheme string) (conf.DataSource, error) {
	creatorFunc, exist := registry[scheme]
	if !exist {
		return nil, ErrInvalidDataSource
	}
	return creatorFunc(), nil
}
