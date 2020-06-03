package defers

import (
	"fmt"
	"sync"

	"github.com/douyu/jupiter/pkg/util/xcolor"
)

var (
	globalDefers = make([]func() error, 0)
	mu           = sync.RWMutex{}
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
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
