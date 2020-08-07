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
	"go.uber.org/multierr"
)

// SerialWithError ...
func SerialWithError(fns ...func() error) func() error {
	return func() error {
		var errs error
		for _, fn := range fns {
			errs = multierr.Append(errs, try(fn, nil))
		}
		return errs
	}
}

// 创建一个迭代器
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn, nil); err != nil {
				return err
				// return errors.Wrap(err, xstring.FunctionName(fn))
			}
		}
		return nil
	}
}

// 策略注入
type WhenError int

var (

	// ReturnWhenError ...
	ReturnWhenError WhenError = 1

	// ContinueWhenError ...
	ContinueWhenError WhenError = 2

	// PanicWhenError ...
	PanicWhenError WhenError = 3

	// LastErrorWhenError ...
	LastErrorWhenError WhenError = 4
)

// SerialWhenError ...
func SerialWhenError(we WhenError) func(fn ...func() error) func() error {
	return func(fns ...func() error) func() error {
		return func() error {
			var errs error
			for _, fn := range fns {
				if err := try(fn, nil); err != nil {
					switch we {
					case ReturnWhenError: // 直接退出
						return err
					case ContinueWhenError: // 继续执行
						errs = multierr.Append(errs, err)
					case PanicWhenError: // panic
						panic(err)
					case LastErrorWhenError: // 返回最后一个错误
						errs = err
					}
				}
			}
			return errs
		}
	}
}
