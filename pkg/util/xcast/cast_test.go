package xcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt64SliceE(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect []int64
		iserr  bool
	}{
		{[]int{1, 3}, []int64{1, 3}, false},
		{[]interface{}{1.2, 3.2}, []int64{1, 3}, false},
		{[]string{"2", "3"}, []int64{2, 3}, false},
		{[2]string{"2", "3"}, []int64{2, 3}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errmsg := fmt.Sprintf("i = %d", i) // assert helper message

		v, err := ToInt64SliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errmsg)
			continue
		}

		assert.NoError(t, err, errmsg)
		assert.Equal(t, test.expect, v, errmsg)

		// Non-E test
		v = ToInt64Slice(test.input)
		assert.Equal(t, test.expect, v, errmsg)
	}
}
