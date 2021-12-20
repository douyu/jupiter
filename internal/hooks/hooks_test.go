package hooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHooksRegister(t *testing.T) {
	var str string
	type args struct {
		fns []func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "register and do",
			args: args{
				fns: []func(){
					func() { str += "1," },
					func() { str += "2," },
					func() { str += "3," },
					func() { str += "4," },
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(Stage_BeforeLoadConfig, tt.args.fns...)
			Do(Stage_BeforeLoadConfig)
			assert.Equal(t, str, "4,3,2,1,")
		})
	}
}

func TestHooksDo(t *testing.T) {
	var str string
	globalHooks[Stage_AfterLoadConfig] = append(globalHooks[Stage_AfterLoadConfig],
		func() { str += "1," },
		func() { str += "2," },
		func() { str += "3," },
		func() { str += "4," },
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
			Do(Stage_AfterLoadConfig)
			Do(7)
			Do(Stage_BeforeStop)
			Do(Stage_AfterStop)
			assert.Equal(t, str, "4,3,2,1,")
		})
	}
}
