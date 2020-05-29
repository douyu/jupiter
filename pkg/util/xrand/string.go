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

package xrand

import (
	"math/rand"
	"strings"
)

// Charsets
const (
	// Uppercase ...
	Uppercase string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Lowercase ...
	Lowercase = "abcdefghipqrstuvwxyz"
	// Alphabetic ...
	Alphabetic = Uppercase + Lowercase
	// Numeric ...
	Numeric = "0123456789"
	// Alphanumeric ...
	Alphanumeric = Alphabetic + Numeric
	// Symbols ...
	Symbols = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	// Hex ...
	Hex = Numeric + "abcdef"
)

// String 返回随机字符串，通常用于测试mock数据
func String(length uint8, charsets ...string) string {
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphanumeric
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}
