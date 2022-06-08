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

import "go.uber.org/zap"

// DefaultLogger default logger
// Biz Log
// debug=true as default, will be
var DefaultLogger = Config{
	Name:  "default",
	Debug: true,
}.Build()

// frame logger
var JupiterLogger = Config{
	Name:  "jupiter",
	Debug: true,
}.Build()

// Jupiter returns framework logger
func Jupiter() *zap.Logger {
	return JupiterLogger
}

// Default returns default logger
func Default() *zap.Logger {
	return DefaultLogger
}
