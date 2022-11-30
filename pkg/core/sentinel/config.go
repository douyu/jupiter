// Copyright 2022 Douyu
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
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/xlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Config struct {
	Enable     bool   `toml:"enable"`
	Datasource string `toml:"datasource"`
	EtcdRawKey string `toml:"etcdRawKey"`
	// 熔断降级
	CbKey   string                `toml:"cbKey"`
	CbRules []*CircuitBreakerRule `toml:"cbRules"`
	// 流量控制
	FlowKey   string       `toml:"flowKey"`
	FlowRules []*flow.Rule `toml:"flowRules"`
	// 系统保护
	SystemKey   string         `toml:"systemKey"`
	SystemRules []*system.Rule `toml:"systemRules"`
}

func StdConfig() Config {
	return RawConfig(constant.ConfigKey("sentinel"))
}

func RawConfig(key string) Config {
	config := DefaultConfig()

	if conf.Get(constant.ConfigKey("sentinel")) == nil {
		return config
	}

	if err := conf.UnmarshalKey(key, &config, conf.TagName("toml")); err != nil {
		xlog.Jupiter().Warn("unmarshal config", zap.String("key", key), zap.Error(err))

		return config
	}

	return config
}

func DefaultConfig() Config {
	return Config{
		Enable:     false,
		Datasource: SENTINEL_DATASOURCE_ETCD,
		EtcdRawKey: "app.registry.etcd",
		// 熔断降级规则 /wsd-sentinel/{language}/{app}/{idc}/{env}/{ruleType}=${value}
		CbKey: "/wsd-sentinel/go/%s/%s/%s/degrade",
		// 流量控制规则 /wsd-sentinel/{language}/{app}/{idc}/{env}/{ruleType}=${value}
		FlowKey: "/wsd-sentinel/go/%s/%s/%s/flow",
		// 系统保护规则 /wsd-sentinel/{language}/{app}/{idc}/{env}/{ruleType}=${value}
		SystemKey: "/wsd-sentinel/go/%s/%s/%s/system",
	}
}

func (e Config) exitHandler(entry *SentinelEntry, ctx *EntryContext) error {
	if ctx.Err() != nil {
		sentinelExceptionsThrown.WithLabelValues(labels(entry.Resource().Name())...).Inc()
	} else {
		sentinelSuccess.WithLabelValues(labels(entry.Resource().Name())...).Inc()
	}

	sentinelRt.WithLabelValues(labels(entry.Resource().Name())...).Observe(float64(ctx.Rt()) / 1000)

	return ctx.Err()
}

func (c Config) Entry(resource string, opts ...EntryOption) (*SentinelEntry, *BlockError) {
	if !c.Enable {
		return base.NewSentinelEntry(nil, nil, nil), nil
	}

	a, b := sentinel.Entry(resource, opts...)

	sentinelReqeust.WithLabelValues(labels(resource)...).Inc()

	if b != nil {
		sentinelBlocked.WithLabelValues(labels(resource)...).Inc()

		return a, b
	}

	a.WhenExit(c.exitHandler)

	return a, b
}

func (c Config) Build() error {

	if !c.Enable {
		xlog.Jupiter().Info("disable sentinel feature")

		return nil
	}

	if err := sentinel.InitDefault(); err != nil {
		xlog.Jupiter().Error("sentinel.InitDefault failed", zap.Error(err))

		return err
	}

	defaultConfig := config.NewDefaultConfig()
	defaultConfig.Sentinel.App.Name = pkg.Name()

	err := sentinel.InitWithConfig(defaultConfig)
	if err != nil {
		return err
	}

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	c.loadRules()

	return nil
}

func (c Config) loadRules() {

	xlog.Jupiter().Info("load sentinel rules", zap.String("datasource", c.Datasource))

	switch c.Datasource {
	case SENTINEL_DATASOURCE_ETCD:

		cli, err := etcdv3.RawConfig(c.EtcdRawKey).Singleton()
		if err != nil {
			panic(err)
		}

		err = initRules(cli.Client, c.CbKey, datasource.NewCircuitBreakerRulesHandler(circuitBreakerRuleJsonArrayParser))
		if err != nil {
			xlog.Jupiter().Warn("sentinel etcd Initialize failed", xlog.FieldErr(err))
		}

		err = initRules(cli.Client, c.SystemKey, datasource.NewSystemRulesHandler(datasource.SystemRuleJsonArrayParser))
		if err != nil {
			xlog.Jupiter().Warn("sentinel etcd Initialize failed", xlog.FieldErr(err))
		}

		err = initRules(cli.Client, c.FlowKey, datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser))
		if err != nil {
			xlog.Jupiter().Warn("sentinel etcd Initialize failed", xlog.FieldErr(err))
		}
	default:

		var err error
		_, err = flow.LoadRules(c.FlowRules)
		if err != nil {
			xlog.Jupiter().Warn("sentinel flow.LoadRules failed", xlog.FieldErr(err))
		}

		_, err = system.LoadRules(c.SystemRules)
		if err != nil {
			xlog.Jupiter().Warn("sentinel system.LoadRules failed", xlog.FieldErr(err))
		}

		rules := convertCbRules(c.CbRules)
		_, err = circuitbreaker.LoadRules(rules)
		if err != nil {
			xlog.Jupiter().Warn("sentinel circuitbreaker.LoadRules failed", xlog.FieldErr(err))
		}
	}
}

func checkSrcComplianceJson(src []byte) (bool, error) {
	if len(src) == 0 {
		return false, nil
	}
	return true, nil
}

func initRules(client *clientv3.Client, key string, h datasource.PropertyHandler) error {
	datasource, err := newDataSource(client,
		fmt.Sprintf(key, pkg.Name(), pkg.AppZone(), conf.GetString("app.mode")),
		h)
	if err != nil {
		return err
	}

	return datasource.Initialize()
}

func circuitBreakerRuleJsonArrayParser(src []byte) (interface{}, error) {
	if valid, err := checkSrcComplianceJson(src); !valid {
		return nil, err
	}

	rules := make([]*CircuitBreakerRule, 0, 8)
	if err := json.Unmarshal(src, &rules); err != nil {
		desc := fmt.Sprintf("Fail to convert source bytes to []*CircuitBreakerRule, err: %s", err.Error())

		xlog.Jupiter().Warn("json.Unmarshal", zap.ByteString("src", src), zap.Error(err))
		return nil, datasource.NewError(datasource.ConvertSourceError, desc)
	}

	xlog.Jupiter().Info("circuitBreakerRuleJsonArrayParser finished", zap.Any("rules", rules))

	return convertCbRules(rules), nil
}
