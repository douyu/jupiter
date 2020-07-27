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

package defers

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xdefer"
)

var (
	globalDefers = xdefer.NewStack()
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
	for _, fn := range fns {
		globalDefers.Push(func() error {
			fmt.Printf("[jupiter] %+v\n", xcolor.Red("clean..."))
			err := fn()
			fmt.Printf("[jupiter] %+v\n", xcolor.Green("clean... done"))
			return err
		})
	}
	globalDefers.Push(fns...)
}

// Clean 清除
func Clean() {
	globalDefers.Clean()
}
