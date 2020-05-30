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
	"fmt"
	"testing"
)

func TestJSON_OmitDefault(t *testing.T) {
	type CC struct {
		D string `json:",omitempty" toml:",omitempty" `
	}
	type AA struct {
		A string `json:"a,omitempty"`
		B int    `json:",omitempty"`
		C CC     `json:",omitempty"`
		D *CC    `json:",omitempty" toml:",omitempty" `
		E *CC    `json:"e" toml:"e" `
	}

	aa := AA{A: "11", D: &CC{}}

	bs, err := OmitDefaultAPI.Marshal(aa)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("string(bs) = %+v\n", string(bs))

	as, err := _jsonAPI.Marshal(aa)
	fmt.Printf("err = %+v\n", err)
	fmt.Printf("as = %+v\n", string(as))
}
