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
	"fmt"

	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

// DebugObject ...
func PrintObject(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		color.RedString(message),
		pretty.Color(
			pretty.Pretty([]byte(xstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}

// DebugBytes ...
func DebugBytes(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(xstring.Json(obj))), pretty.TerminalStyle))
}

// PrintKV ...
func PrintKV(key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-50s => %s\n", color.RedString(key), color.GreenString(val))
}

// PrettyKVWithPrefix ...
func PrintKVWithPrefix(prefix string, key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-8s]> %-30s => %s\n", prefix, color.RedString(key), color.BlueString(val))
}

// PrintMap ...
func PrintMap(data map[string]interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	for key, val := range data {
		fmt.Printf("%-20s : %s\n", color.RedString(key), fmt.Sprintf("%+v", val))
	}
}
