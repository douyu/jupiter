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

import (
	"errors"
)

// Attributes ...
type Attributes struct {
	m map[interface{}]interface{}
}

var (
	// ErrInvalidKVPairs ...
	ErrInvalidKVPairs = errors.New("invalid kv pairs")
)

// New ...
func New(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	a := &Attributes{m: make(map[interface{}]interface{}, len(kvs)/2)}
	for i := 0; i < len(kvs)/2; i++ {
		a.m[kvs[i*2]] = kvs[i*2+1]
	}
	return a
}

// WithValues ...
func (a *Attributes) WithValues(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	n := &Attributes{m: make(map[interface{}]interface{}, len(a.m)+len(kvs)/2)}
	for k, v := range a.m {
		n.m[k] = v
	}
	for i := 0; i < len(kvs)/2; i++ {
		n.m[kvs[i*2]] = kvs[i*2+1]
	}
	return n
}

// Value ...
func (a *Attributes) Value(key interface{}) interface{} {
	return a.m[key]
}
