package xcolor

import "fmt"

// array transform
func arrToTransform(arg []interface{}) interface{} {
	var res interface{}

	for _, v := range arg {
		if res != nil {
			res = fmt.Sprintf("%v %v", res, v)
		} else {
			res = v
		}
	}

	return res
}
