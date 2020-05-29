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

package xcast

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrUnaddressable unaddressable val
	ErrUnaddressable = errors.New("val must be addressable")
	// ErrNotPointer pinter val
	ErrNotPointer = errors.New("val must be a pointer")
)

// Decode decode interface into struct
func Decode(m interface{}, val interface{}) error {
	if err := check(val); err != nil {
		return err
	}

	return decode(m, reflect.ValueOf(val).Elem())
}

func check(val interface{}) error {
	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	if !v.Elem().CanAddr() {
		return ErrUnaddressable
	}

	return nil
}

func decode(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	if !dataVal.IsValid() {
		val.Set(reflect.Zero(dataVal.Type()))
		return nil
	}

	kind := val.Kind()
	switch kind {
	case reflect.Bool:
		return decodeBool(data, val)
	case reflect.Interface:
		return decodeInterface(data, val)
	case reflect.String:
		return decodeString(data, val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decodeInt(data, val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decodeUint(data, val)
	case reflect.Float32, reflect.Float64:
		return decodeFloat(data, val)
	case reflect.Map, reflect.Slice:
		return decodeInterface(data, val)
	case reflect.Ptr:
		return decodePtr(data, val)
	case reflect.Struct:
		return decodeStruct(data, val)
	default:
		return fmt.Errorf("unsupported type %s", kind)
	}
}

func decodeStruct(data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := dataVal.Kind()

	if dataVal.Type() == val.Type() {
		val.Set(dataVal)
		return nil
	}

	switch dataKind {
	// Only map can converted into struct
	case reflect.Map:

	default:
		return fmt.Errorf("")
	}

	return nil
}

func decodePtr(data interface{}, val reflect.Value) error {
	valType := val.Type()
	valElem := valType.Elem()

	elem := val
	if elem.IsNil() {
		elem = reflect.New(valElem)
	}

	if err := decode(data, reflect.Indirect(elem)); err != nil {
		return err
	}

	val.Set(elem)
	return nil
}

func decodeBool(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Bool:
		val.SetBool(dataVal.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetBool(0 != dataVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val.SetBool(0 != dataVal.Uint())
	case reflect.Float32, reflect.Float64:
		val.SetBool(0 != dataVal.Float())
	case reflect.String:
		ok := strings.Contains(" True true ", dataVal.String())
		val.SetBool(ok)
	default:
		return fmt.Errorf("")

	}
	return nil

}

func decodeInt(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Bool:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetInt(dataVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val.SetInt(int64(dataVal.Uint()))
	case reflect.Float32, reflect.Float64:
		val.SetInt(int64(dataVal.Float()))
	case reflect.String:
		d, err := strconv.ParseInt(dataVal.String(), 0, val.Type().Bits())
		if err != nil {
			return fmt.Errorf("parse '%s' as int failed: %s", dataVal.String(), err)
		}
		val.SetInt(d)
	default:
		return fmt.Errorf("decode int failed: %#v", data)
	}
	return nil
}

func decodeUint(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Bool:
		if dataVal.Bool() {
			val.SetUint(1)
		} else {
			val.SetUint(0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := dataVal.Int()
		if i < 0 {
			return fmt.Errorf("decode uint failed: int out of range '%d'", i)
		}
		val.SetUint(uint64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val.SetUint(dataVal.Uint())
	case reflect.Float32, reflect.Float64:
		f := dataVal.Float()
		if f < 0 {
			return fmt.Errorf("decode uint failed: int out of range '%f'", f)
		}
		val.SetUint(uint64(f))
	case reflect.String:
		d, err := strconv.ParseUint(dataVal.String(), 0, val.Type().Bits())
		if err != nil {
			return fmt.Errorf("parse '%s' as int failed: %s", dataVal.String(), err)
		}
		val.SetUint(d)
	default:
		return fmt.Errorf("decode uint failed: %#v", data)
	}
	return nil
}

func decodeString(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Bool:
		if dataVal.Bool() {
			val.SetString("1")
		} else {
			val.SetString("0")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	case reflect.String:
		d, err := strconv.ParseInt(dataVal.String(), 0, val.Type().Bits())
		if err != nil {
			return fmt.Errorf("parse '%s' as int failed: %s", dataVal.String(), err)
		}
		val.SetInt(d)
	default:
		return fmt.Errorf("decode int failed: %#v", data)
	}
	return nil
}

func decodeInterface(data interface{}, val reflect.Value) error {
	valType := val.Type()
	valKey := valType.Key()
	valElem := valType.Elem()

	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Map:
		valMap := val
		if valMap.IsNil() {
			valMap = reflect.MakeMap(reflect.MapOf(valKey, valElem))
		}
		for _, k := range dataVal.MapKeys() {
			subKey := reflect.Indirect(reflect.New(valType))
			if err := decode(k.Interface(), subKey); err != nil {
				continue
			}

			v := dataVal.MapIndex(k).Interface()
			subVal := reflect.Indirect(reflect.New(valElem))
			if err := decode(v, subVal); err != nil {
				continue
			}

			valMap.SetMapIndex(subKey, subVal)
		}

		val.Set(valMap)
	case reflect.Array:
	case reflect.Slice:
		valSlice := val
		if valSlice.IsNil() {
			valSlice = reflect.MakeSlice(reflect.SliceOf(valElem), dataVal.Len(), dataVal.Len())

		}
		for i := 0; i < dataVal.Len(); i++ {
			subData := dataVal.Index(i).Interface()
			for valSlice.Len() <= i {
				valSlice = reflect.Append(valSlice, reflect.Zero(valElem))
			}
			subField := valSlice.Index(i)
			if err := decode(subData, subField); err != nil {
				continue
			}
		}

		val.Set(valSlice)
	default:
		return fmt.Errorf("decode map failed: %#v", data)
	}

	return nil
}

func decodeFloat(data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := dataVal.Kind()
	switch dataKind {
	case reflect.Bool:
		if dataVal.Bool() {
			val.SetFloat(1.0)
		} else {
			val.SetFloat(0.0)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetFloat(float64(dataVal.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val.SetFloat(float64(dataVal.Uint()))
	case reflect.Float32, reflect.Float64:
		val.SetFloat(dataVal.Float())
	case reflect.String:
		f, err := strconv.ParseFloat(dataVal.String(), val.Type().Bits())
		if err != nil {
			return fmt.Errorf("parse '%s' as int failed: %s", dataVal.String(), err)
		}
		val.SetFloat(f)
	default:
		return fmt.Errorf("decode int failed: %#v", data)
	}
	return nil
}
