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

package xattr

import "testing"

// New ...
func TestNew(t *testing.T) {
	k1 := 1
	v1 := "first"
	attr := New(k1, v1)

	ret1, ok1 := attr.Value(k1).(string)
	if !ok1 || v1 != ret1 {
		t.Fatalf("attr.Value error: want:%v ret:%v", v1, ret1)
	}

	k2 := "2"
	v2 := 2
	attr = attr.WithValues(k2, v2)
	ret2, ok2 := attr.Value(k2).(int)
	if !ok2 || v2 != ret2 {
		t.Fatalf("attr.WithValues error: want:%v ret:%v", v2, ret2)
	}
}
