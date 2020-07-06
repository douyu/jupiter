package defers

import (
	"testing"

	"github.com/douyu/jupiter/pkg/util/xdefer"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	var str string
	type args struct {
		fns []func() error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "register",
			args: args{
				fns: []func() error{
					func() error { str += "1,"; return nil },
					func() error { str += "2,"; return nil },
					func() error { str += "3,"; return nil },
					func() error { str += "4,"; return nil },
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.fns...)
		})
	}
}

func TestClean(t *testing.T) {
	var str string
	globalDefers = xdefer.NewStack()
	globalDefers.Push(
		func() error { str += "1,"; return nil },
		func() error { str += "2,"; return nil },
		func() error { str += "3,"; return nil },
		func() error { str += "4,"; return nil },
	)

	tests := []struct {
		name string
	}{
		{
			"testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Clean()
			assert.Equal(t, str, "4,3,2,1,")
		})
	}
}
