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
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64
	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8
	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8

	// NodeID determine the specific value of the node id; 0 nodeID are not used by default
	NodeID int64

	// EnableRedis whether to enable the use of redis to assign node IDs
	EnableRedis bool
	// RedisAddr redis addr, default to 'Host:Port'
	RedisAddr string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Epoch:    snowflake.Epoch,
		NodeBits: snowflake.NodeBits,
		StepBits: snowflake.StepBits,
		NodeID:   flag.Int("nodeId"),
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
	if config.Epoch != 0 {
		snowflake.Epoch = config.Epoch
	}

	if config.NodeBits != 0 {
		snowflake.NodeBits = config.NodeBits
	}

	if config.StepBits != 0 {
		snowflake.StepBits = config.StepBits
	}

	if snowflake.NodeBits+snowflake.StepBits != 22 {
		return nil, fmt.Errorf("snowflake NodeBits:%v and StepBits:%v err,sum dont 22", snowflake.NodeBits, snowflake.StepBits)
	}

	if config.NodeID == 0 {
		// use the default node id -> 1
		config.NodeID = 1
	}

	return &Uuid{
		snowflakeRw: &sync.RWMutex{},
		nodeId:      config.NodeID,
		enableRedis: config.EnableRedis,
	}, nil
}
