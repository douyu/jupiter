package xdefer

import (
	"testing"
)

func TestNewStack(t *testing.T) {
	// stack := &DeferStack{
	// 	fns: make([]func() error, 0),
	// 	mu:  sync.RWMutex{},
	// }
	stack := NewStack()
	state := ""
	fn1 := func() error {
		state = state + "1"
		return nil
	}
	fn2 := func() error {
		state = state + "2"
		return nil
	}
	fn3 := func() error {
		state = state + "3"
		return nil
	}
	stack.Push(fn1, fn2)
	stack.Push(fn3)
	stack.Clean()
	want := "321"
	if state != want {
		t.Fatalf("Stack has error,want:%v ret:%v", want, state)
	}
}
