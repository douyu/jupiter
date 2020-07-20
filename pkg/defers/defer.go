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
