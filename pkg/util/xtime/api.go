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

package xtime

import "time"

// GetTimestampInMilli ...
func GetTimestampInMilli() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}

// Timing the cost of function call, unix nano was returned
func Elapse(f func()) int64 {
	now := time.Now().UnixNano()
	f()
	return time.Now().UnixNano() - now
}

// IsLeapYear ...
func IsLeapYear(year int) bool {
	if year%100 == 0 {
		return year%400 == 0
	}

	return year%4 == 0
}
