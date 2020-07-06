package xwaitgroup

import "sync"

type XWaitGroup struct {
	sync.WaitGroup
}

func (w *XWaitGroup) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
