package xerr

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
)

// ErrGoexit indicates the runtime.Goexit was called in
// the user given function.
var ErrGoexit = errors.New("runtime.Goexit was called")

// A PanicError is an arbitrary value recovered from a panic
// with the stack trace during the execution of given function.
type PanicError struct {
	value interface{}
	stack []byte
}

// Error implements error interface.
func (p *PanicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func NewPanicError(v interface{}) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &PanicError{value: v, stack: stack}
}
