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

package xlog

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func BenchmarkWriteSyncer(b *testing.B) {
	b.Run("write file with no buffer", func(b *testing.B) {
		file, err := ioutil.TempFile("", "log")
		assert.NoError(b, err)
		defer file.Close()
		defer os.Remove(file.Name())

		w := zapcore.AddSync(file)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
	b.Run("write file with buffer", func(b *testing.B) {
		file, err := ioutil.TempFile("", "log")
		assert.NoError(b, err)
		defer file.Close()
		defer os.Remove(file.Name())

		w, close := Buffer(zapcore.AddSync(file), 0, 0)
		defer close()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
}
