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

package xstruct

import "reflect"

// CopyStruct ...
func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		dstValue := dstVal.FieldByName(name)
		if dstValue.IsValid() == false {
			continue
		}
		dstValue.Set(value)
	}
}
