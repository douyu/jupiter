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

package xstring

import (
	"strings"
	"unicode/utf8"
)

// ToSnakeCase 转为snake格式
func ToSnakeCase(str string) string {
	str = strings.TrimSpace(strings.ToLower(str))
	return strings.Replace(str, " ", "_", -1)
}

// ToCamelCase 转为驼峰格式
func ToCamelCase(str string) string {
	str = strings.TrimSpace(str)
	if utf8.RuneCountInString(str) < 2 {
		return str
	}

	var buff strings.Builder
	var temp string
	for _, r := range str {
		c := string(r)
		if c != " " {
			if temp == " " {
				c = strings.ToUpper(c)
			}
			_, _ = buff.WriteString(c)
		}
		temp = c
	}
	return buff.String()
}
