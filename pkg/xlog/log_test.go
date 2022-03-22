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

package xlog_test

import (
	"context"
	"testing"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/stretchr/testify/assert"
)

func Test_log(t *testing.T) {
	xlog.Debug("debug", xlog.Any("a", "b"))
	xlog.Info("info", xlog.Any("a", "b"))
	xlog.Warn("warn", xlog.Any("a", "b"))
	xlog.Error("error", xlog.Any("a", "b"))

	xlog.Debugw("debug", xlog.Any("a", "b"))
	xlog.Infow("info", xlog.Any("a", "b"))
	xlog.Warnw("warn", xlog.Any("a", "b"))
	xlog.Errorw("error", xlog.Any("a", "b"))

	xlog.Debugf("debug", xlog.Any("a", "b"))
	xlog.Infof("info", xlog.Any("a", "b"))
	xlog.Warnf("warn", xlog.Any("a", "b"))
	xlog.Errorf("error", xlog.Any("a", "b"))
}

func Test_trace(t *testing.T) {
	ctx := xlog.SetTraceID(context.TODO(), "traceid")
	assert.Equal(t, xlog.GetTraceID(ctx), "traceid")
	xlog.FromContext(ctx).Debug("debug", xlog.Any("a", "b"))
}
