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

package xcast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	A string `json:"a"`
	B string `json:"b"`
	C struct {
		D []string `json:"d"`
	} `json:"c"`
}

func Test_Decode(t *testing.T) {
	var src = map[string]interface{}{
		"a": "1",
		"b": 2,
		"c": map[string]interface{}{
			"d": []string{"1", "2", "3"},
		},
	}

	var p T

	err := Decode(src, &p)
	assert.Nil(t, err)
}
