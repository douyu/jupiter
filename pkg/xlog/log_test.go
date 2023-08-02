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

package xlog

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func Test_log(t *testing.T) {

	stdlog := Default()
	stdlog.Debug("debug", Any("a", "b"))
	stdlog.Info("info", Any("a", "b"))
	stdlog.Warn("warn", Any("a", "b"))
	stdlog.Error("error", Any("a", "b"))

	data, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(t, err)

	found := false
	for _, v := range data {

		if v.GetName() == "log_level_total" {
			assert.NotEmpty(t, v.GetMetric())
			found = true

			fmt.Println(v.GetMetric())
			for _, vv := range v.GetMetric() {
				if vv.Counter.GetValue() != 1 {
					t.Fail()
				}
			}
		}
	}

	// no metrics found
	if !found {
		assert.FailNow(t, "should never reach here")
	}
}

func Test_trace(t *testing.T) {

	log := Jupiter()
	ctx := NewContext(context.TODO(), log, "a:b:c:1")

	jlog := J(ctx)
	jlog.Debug("debug", Any("a", "b"))
	jlog.Info("info", Any("a", "b"))
	jlog.Warn("warn", Any("a", "b"))
	jlog.Error("error", Any("a", "b"))

	data, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(t, err)

	found := false
	for _, v := range data {

		if v.GetName() == "log_level_total" {
			assert.NotEmpty(t, v.GetMetric())
			found = true

			fmt.Println(v.GetMetric())
			for _, vv := range v.GetMetric() {
				if vv.Counter.GetValue() != 1 {
					assert.FailNow(t, "should be 1")
				}
			}
		}
	}

	// no metrics found
	if !found {
		assert.FailNow(t, "should never reach here")
	}
}

func TestXlog(t *testing.T) {
	defaultConfig := DefaultConfig()
	core, olog := observer.New(zapcore.InfoLevel)
	defaultConfig.Core = core
	log := defaultConfig.Build()

	log.Debug("debug", Any("a", "b"))
	log.Info("info", Any("a", "b"), FieldCost(time.Second))
	log.Warn("warn", Any("a", "b"))
	log.Error("error", Any("a", "b"))

	assert.Equal(t, 3, len(olog.All()))
	assert.Equal(t, "info", olog.All()[0].Message)
	assert.Equal(t, "b", olog.All()[0].ContextMap()["a"])
	assert.Equal(t, "1000.000", olog.All()[0].ContextMap()["cost"])
	assert.Equal(t, "1234567890", olog.All()[0].ContextMap()["aid"])
}
