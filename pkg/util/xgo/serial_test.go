// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xgo

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestSerialE(t *testing.T) {
	type args struct {
		fns []func() error
	}
	fn := func() error {
		return errors.New("error")
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "serial",
			args: args{
				fns: []func() error{
					fn, fn, fn, fn,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SerialWithError(tt.args.fns...)
			assert.NotNil(t, got)
			err := got()
			assert.NotNil(t, err)
			errs := multierr.Errors(err)
			assert.Len(t, errs, 4)
			for _, err := range errs {
				assert.Equal(t, err.Error(), "error")
			}
		})
	}
}

func TestSerialUntilError(t *testing.T) {
	type args struct {
		fns []func() error
	}
	var value int64
	fn := func(arg int64) func() error {
		return func() error {
			if arg < 0 {
				return errors.New("invalid")
			}
			atomic.AddInt64(&value, arg)
			return nil
		}
	}
	tests := []struct {
		name string
		args args
		want func() error
	}{
		{
			name: "until",
			args: args{
				fns: []func() error{
					fn(1), fn(2), fn(-1), fn(4),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SerialUntilError(tt.args.fns...)
			err := got()
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), "invalid")
			assert.Equal(t, atomic.LoadInt64(&value), int64(1+2))
		})
	}
}

func TestSerialWhenError(t *testing.T) {
	type args struct {
		fns []func() error
	}
	var value int64
	fn := func(arg int64) func() error {
		return func() error {
			if arg < 0 {
				return fmt.Errorf("invalid %+d", arg)
			}
			atomic.AddInt64(&value, arg)
			return nil
		}
	}
	tests := []struct {
		name string
		args args
		we   WhenError
	}{
		// TODO: Add test cases.
		{
			name: "panic when error",
			args: args{
				fns: []func() error{
					fn(1), fn(-1), fn(3), fn(-2), fn(5),
				},
			},
			we: PanicWhenError,
		},
		{
			name: "continue when error",
			args: args{
				fns: []func() error{
					fn(1), fn(-1), fn(3), fn(-2), fn(5),
				},
			},
			we: ContinueWhenError,
		},
		{
			name: "return when error",
			args: args{
				fns: []func() error{
					fn(1), fn(-1), fn(3), fn(-2), fn(5),
				},
			},
			we: ReturnWhenError,
		},
		{
			name: "return last err when error",
			args: args{
				fns: []func() error{
					fn(1), fn(-1), fn(3), fn(-2), fn(5),
				},
			},
			we: LastErrorWhenError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.StoreInt64(&value, 0)
			got := SerialWhenError(tt.we)
			switch tt.we {
			case ContinueWhenError:
				err := got(tt.args.fns...)()
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "invalid -1; invalid -2")
				assert.Equal(t, atomic.LoadInt64(&value), int64(1+3+5))
			case PanicWhenError:
				assert.Panics(t, func() { _ = got(tt.args.fns...)() })
			case LastErrorWhenError:
				err := got(tt.args.fns...)()
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "invalid -2")
			case ReturnWhenError:
				err := got(tt.args.fns...)()
				assert.NotNil(t, err)
				assert.Equal(t, atomic.LoadInt64(&value), int64(1))
			default:
				t.Fail()
			}
		})
	}
}
