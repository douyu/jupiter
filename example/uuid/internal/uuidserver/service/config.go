package service

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/pkg/errors"
)

// ModName named a mod
const ModName = "uuid"

func init() {
	flag.Register(&flag.IntFlag{
		Name:    "nodeId",
		Usage:   "--nodeId, set uuid service nodeId",
		Default: 1,
		Action: func(name string, fs *flag.FlagSet) {
			xlog.Debugf("nodeId flag: %v", fs.Int(name))
		},
	})
}

// Config uuid service config
type Config struct {
	NodeId int64
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		NodeId: flag.Int("nodeId"),
	}
}

// StdConfig Jupiter Standard HTTP Server config
func StdConfig(name string) *Config {
	return RawConfig("jupiter.server." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil &&
		errors.Cause(err) != conf.ErrInvalidKey {
		xlog.Panic("http server parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

func (config *Config) MustBuild() *Uuid {
	server, err := config.Build()
	if err != nil {
		xlog.Panicf("build uuid server failed: %v", err)
	}
	return server
}

// Build create server instance, then initialize it with necessary interceptor
func (config *Config) Build() (*Uuid, error) {
	// 判断 NodeID 的值 是否正常

	// Create a new Node with a Node number of nodeId
	node, err := snowflake.NewNode(config.NodeId)
	if err != nil {
		return nil, fmt.Errorf("snowflake NewNode err:%v", err)
	}

	return &Uuid{
		snowflakeRw:  &sync.RWMutex{},
		snowflakeMap: node,
	}, nil
}
