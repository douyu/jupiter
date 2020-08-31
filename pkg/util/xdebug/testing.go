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

package xdebug

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/tidwall/pretty"

	"github.com/douyu/jupiter/pkg/xlog"
)

var (
	isTestingMode     bool
	isDevelopmentMode = os.Getenv("JUPITER_MODE") == "dev"
)

func init() {
	if isDevelopmentMode {
		xlog.DefaultLogger.SetLevel(xlog.DebugLevel)
		xlog.JupiterLogger.SetLevel(xlog.DebugLevel)
	}
}

// IsTestingMode 判断是否在测试模式下
var onceTest = sync.Once{}

// IsTestingMode ...
func IsTestingMode() bool {
	onceTest.Do(func() {
		isTestingMode = flag.Lookup("test.v") != nil
	})

	return isTestingMode
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return isDevelopmentMode || isTestingMode
}

// IfPanic ...
func IfPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyJsonPrint ...
func PrettyJsonPrint(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		xcolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(xstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}
