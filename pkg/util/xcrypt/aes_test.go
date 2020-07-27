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

package xcrypt

import "testing"

const key = "123456781234567812345678"

func TestAesEncrypt(t *testing.T) {
	pwd := "jupiter"
	t.Run("jupiter", func(t *testing.T) {
		t.Log(AesEncrypt([]byte(pwd), key))
	})

}

func BenchmarkAesEncrypt(b *testing.B) {
	pwd := "jupiter"
	b.Log(AesEncrypt([]byte(pwd), key))
}

func TestAesDecrypt(t *testing.T) {
	aesPwd := "94knaB7YbCkh7S41uS9qNQ"
	t.Log(AesDecrypt(aesPwd, key))
}

func BenchmarkAesDecrypt(b *testing.B) {
	aesPwd := "94knaB7YbCkh7S41uS9qNQ"
	b.Log(AesDecrypt(aesPwd, key))
}
