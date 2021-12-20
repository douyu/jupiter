package hooks

import (
	"fmt"
	"sync"

	"github.com/douyu/jupiter/pkg/util/xcolor"
)

var (
	globalHooks = make(map[Stage][]func())
	mu          = sync.RWMutex{}
)

type Stage int

func (s Stage) String() string {
	switch s {
	case Stage_BeforeLoadConfig:
		return "BeforeLoadConfig"
	case Stage_AfterLoadConfig:
		return "AfterLoadStart"
	case Stage_BeforeStop:
		return "BeforeStop"
	case Stage_AfterStop:
		return "AfterStop"
	}

	return "Unkown"
}

const (
	Stage_BeforeLoadConfig Stage = iota + 1
	Stage_AfterLoadConfig
	Stage_BeforeStop
	Stage_AfterStop
)

// Register 注册一个defer函数
func Register(stage Stage, fns ...func()) {
	mu.Lock()
	defer mu.Unlock()

	globalHooks[stage] = append(globalHooks[stage], fns...)
}

// Do 执行
func Do(stage Stage) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("[jupiter] %+v\n", xcolor.Green(fmt.Sprintf("hook stage (%s)...", stage)))

	for i := len(globalHooks[stage]) - 1; i >= 0; i-- {
		fn := globalHooks[stage][i]
		if fn != nil {
			fn()
		}
	}

	fmt.Printf("[jupiter] %+v\n", xcolor.Green(fmt.Sprintf("hook stage (%s)... done", stage)))
}
