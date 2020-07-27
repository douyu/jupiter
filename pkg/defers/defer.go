package defers

import (
	"github.com/douyu/jupiter/pkg/util/xdefer"
)

var (
	globalDefers = xdefer.NewStack()
)

// Register 注册一个defer函数
func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}

// Clean 清除
func Clean() {
	globalDefers.Clean()
}
