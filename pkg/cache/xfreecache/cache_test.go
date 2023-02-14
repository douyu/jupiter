package xfreecache

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Age  int
	Name string
}

func TestLocalCache(t *testing.T) {
	oneCache := DefaultConfig().Build()
	missCount := 0

	tests := []struct {
		stu Student
	}{
		{
			stu: Student{
				Age:  1,
				Name: "Student 1",
			},
		},
		{
			stu: Student{
				Age:  2,
				Name: "Student 2",
			},
		},
		{
			stu: Student{
				Age:  1,
				Name: "Student 1",
			},
		},
		{
			stu: Student{
				Age:  1,
				Name: "Student 3",
			},
		},
		{
			stu: Student{
				Age:  2,
				Name: "Student 2",
			},
		},
	}

	for _, tt := range tests {
		key := fmt.Sprintf("%d-%s", tt.stu.Age, tt.stu.Name)
		result, _ := oneCache.GetAndSetCacheData(key, func() ([]byte, error) {
			missCount++
			fmt.Println("local cache miss hit")
			ret, _ := json.Marshal(tt.stu)
			return ret, nil
		})
		ret := Student{}
		_ = json.Unmarshal(result, &ret)
		fmt.Println(ret)
		assert.Equalf(t, tt.stu, ret, "GetAndSetCacheData(%v) cache value error", key)
	}
	assert.Equalf(t, missCount, 3, "GetAndSetCacheData miss count error")
}

// BenchmarkLocalCache 1553 ns/op
func BenchmarkLocalCache(b *testing.B) {
	stu := Student{
		Age:  1,
		Name: "stu1",
	}
	oneCache := DefaultConfig().Build()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d-%s", stu.Age, stu.Name)
		result, _ := oneCache.GetAndSetCacheData(key, func() ([]byte, error) {
			ret, _ := json.Marshal(stu)
			return ret, nil
		})
		ret := Student{}
		_ = json.Unmarshal(result, &ret)
	}
}
