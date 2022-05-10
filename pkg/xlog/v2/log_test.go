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

package xlog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/douyu/jupiter/pkg/xlog/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func Test_log(t *testing.T) {

	stdlog := xlog.DefaultLogger
	stdlog.Debug("debug", xlog.Any("a", "b"))
	stdlog.Info("info", xlog.Any("a", "b"))
	stdlog.Warn("warn", xlog.Any("a", "b"))
	stdlog.Error("error", xlog.Any("a", "b"))

	data, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(t, err)

	found := false
	for _, v := range data {

		if v.GetName() == "jupiter_log_level_total" {
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

	log := xlog.DefaultLogger.With(xlog.String("traceid", "a:b:c:1"))
	ctx := xlog.NewContext(context.TODO(), log)

	stdlog := xlog.FromContext(ctx)
	stdlog.Debug("debug", xlog.Any("a", "b"))
	stdlog.Info("info", xlog.Any("a", "b"))
	stdlog.Warn("warn", xlog.Any("a", "b"))
	stdlog.Error("error", xlog.Any("a", "b"))

	data, err := prometheus.DefaultGatherer.Gather()
	assert.Nil(t, err)

	found := false
	for _, v := range data {

		if v.GetName() == "jupiter_log_level_total" {
			assert.NotEmpty(t, v.GetMetric())
			found = true

			fmt.Println(v.GetMetric())
			for _, vv := range v.GetMetric() {
				if vv.Counter.GetValue() != 2 {
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
