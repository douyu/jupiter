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

// +build darwin

package xcolor

import (
	"fmt"
	"math/rand"
	"strconv"
)

var _ = RandomColor()

// RandomColor generates a random color.
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow ...
func Yellow(msg string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", msg)
}

// Red ...
func Red(msg string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", msg)
}

// Redf ...
func Redf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m %+v\n", msg, arg)
}

// Blue ...
func Blue(msg string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", msg)
}

// Green ...
func Green(msg string) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", msg)
}

// Greenf ...
func Greenf(msg string, arg interface{}) string {
	return fmt.Sprintf("\x1b[32m%s\x1b[0m %+v\n", msg, arg)
}
