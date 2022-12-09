package xcast

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// ToInt64Slice casts an interface to a []int64 type.
func ToInt64Slice(i interface{}) []int64 {
	v, _ := ToInt64SliceE(i)
	return v
}

// ToInt64SliceE casts an empty interface to a []int64.
func ToInt64SliceE(i interface{}) ([]int64, error) {
	if i == nil {
		return []int64{}, fmt.Errorf("Unable to Cast %#v to []int64", i)
	}

	switch v := i.(type) {
	case []int64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt64E(s.Index(j).Interface())
			if err != nil {
				return []int64{}, fmt.Errorf("Unable to Cast %#v to []int64", i)
			}
			a[j] = val

		}
		return a, nil
	default:
		return []int64{}, fmt.Errorf("Unable to Cast %#v to []int64", i)
	}
}
