package governor

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/xlog"
)

// ModName ..
const ModName = "govern"

// Config ...
type Config struct {
	Host    string
	Port    int
	Network string `json:"network" toml:"network"`
	logger  *xlog.Logger
	Enable  bool

	// ServiceAddress service address in registry info, default to 'Host:Port'
	ServiceAddress string
}

// StdConfig represents Standard gRPC Server config
// which will parse config by conf package,
// panic if no config key found in conf
func StdConfig(name string) *Config {
	return RawConfig(constant.ConfigKey("server." + name))
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Get(key) == nil {
		return config
	}
	if err := conf.UnmarshalKey(key, &config); err != nil {
		config.logger.Panic("govern server parse config panic",
			xlog.FieldErr(err), xlog.FieldKey(key),
			xlog.FieldValueAny(config),
		)
	}
	return config
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {

	return &Config{
		Enable:  true,
		Host:    flag.String("host"),
		Network: "tcp4",
		Port:    9093,
		logger:  xlog.Jupiter().With(xlog.FieldMod(ModName)),
	}
}

// Build ...
func (config *Config) Build() *Server {
	return newServer(config)
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
