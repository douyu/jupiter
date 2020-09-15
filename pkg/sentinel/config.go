// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentinel

import (
	"encoding/json"
	"io/ioutil"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	sentinel_config "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
)

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("jupiter.sentinel." + name)
}

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
func (config *Config) Build() error {
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

	configEntity := sentinel_config.NewDefaultConfig()
	configEntity.Sentinel.App.Name = config.AppName
	configEntity.Sentinel.Log.Dir = config.LogPath

	if len(config.FlowRules) > 0 {
		_, _ = flow.LoadRules(config.FlowRules)
	}
	return sentinel.InitWithConfig(configEntity)
}

func Entry(resource string) (*base.SentinelEntry, *base.BlockError) {
	return sentinel.Entry(resource)
}
