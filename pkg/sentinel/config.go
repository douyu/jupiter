package sentinel

import (
	"encoding/json"
	"io/ioutil"
	"os"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/xlog"
)

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		xlog.Panic("unmarshal key", xlog.Any("err", err))
	}
	return config
}

// Config ...
type Config struct {
	AppName       string           `json:"appName"`
	LogPath       string           `json:"logPath"`
	FlowRules     []*flow.FlowRule `json:"rules"`
	FlowRulesFile string           `json:"flowRulesFile"`
}

// DefaultConfig returns default config for sentinel
func DefaultConfig() *Config {
	return &Config{
		AppName:   pkg.Name(),
		LogPath:   "/tmp/log",
		FlowRules: make([]*flow.FlowRule, 0),
	}
}

// InitSentinelCoreComponent init sentinel core component
// Currently, only flow rules from json file is supported
// todo: support dynamic rule config
// todo: support more rule such as system rule
func (config *Config) InitSentinelCoreComponent() error {
	if config.FlowRulesFile != "" {
		var rules []*flow.FlowRule
		content, err := ioutil.ReadFile(config.FlowRulesFile)
		if err != nil {
			xlog.Error("load sentinel flow rules", xlog.FieldErr(err), xlog.FieldKey(config.FlowRulesFile))
		}

		if err := json.Unmarshal(content, &rules); err != nil {
			xlog.Error("load sentinel flow rules", xlog.FieldErr(err), xlog.FieldKey(config.FlowRulesFile))
		}

		config.FlowRules = append(config.FlowRules, rules...)
	}

	_ = os.Setenv(constant.EnvKeySentinelAppName, config.AppName)
	_ = os.Setenv(constant.EnvKeySentinelLogDir, config.LogPath)

	if len(config.FlowRules) > 0 {
		_, _ = flow.LoadRules(config.FlowRules)
	}

	return sentinel.InitDefault()
}
