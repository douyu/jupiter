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
	"sync"

	"github.com/douyu/jupiter/pkg/util/xcolor"
)

var (
	globalDefers = make([]func(), 0)
	mu           = sync.RWMutex{}
)

// Register 注册一个defer函数
func Register(fns ...func()) {
	mu.Lock()
	defer mu.Unlock()

	globalDefers = append(globalDefers, fns...)
}

// Clean 清除
func Clean() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("[jupiter] %+v\n", xcolor.Red("clean..."))

	for i := len(globalDefers) - 1; i >= 0; i-- {
		globalDefers[i]()
	}
	fmt.Printf("[jupiter] %+v\n", xcolor.Green("clean... done"))
}
