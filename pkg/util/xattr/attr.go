package xattr

import (
	"errors"
)

// Attributes ...
type Attributes struct {
	m map[interface{}]interface{}
}

var (
	// ErrInvalidKVPairs ...
	ErrInvalidKVPairs = errors.New("invalid kv pairs")
)

// New ...
func New(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	a := &Attributes{m: make(map[interface{}]interface{}, len(kvs)/2)}
	for i := 0; i < len(kvs)/2; i++ {
		a.m[kvs[i*2]] = kvs[i*2+1]
	}
	return a
}

// WithValues ...
func (a *Attributes) WithValues(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	n := &Attributes{m: make(map[interface{}]interface{}, len(a.m)+len(kvs)/2)}
	for k, v := range a.m {
		n.m[k] = v
	}
	for i := 0; i < len(kvs)/2; i++ {
		n.m[kvs[i*2]] = kvs[i*2+1]
	}
	return n
}

// Value ...
func (a *Attributes) Value(key interface{}) interface{} {
	return a.m[key]
}
