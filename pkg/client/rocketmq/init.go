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

package rocketmq

import (
	"fmt"
	"os"
	"runtime"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
)

func init() {

	primitive.PanicHandler = func(i interface{}) {
		stack := make([]byte, 1024)
		length := runtime.Stack(stack, true)
		fmt.Fprint(os.Stderr, "[rocketmq panic recovery]\n", string(stack[:length]))
		xlog.Jupiter().Named(ecode.ModeClientRocketMQ).Error("rocketmq panic recovery", zap.Any("error", i))
	}
}
